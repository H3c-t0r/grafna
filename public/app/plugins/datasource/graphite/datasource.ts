import _ from 'lodash';
import { dateMath, ScopedVars } from '@grafana/data';
import { isVersionGtOrEq, SemVersion } from 'app/core/utils/version';
import gfunc from './gfunc';
import { IQService } from 'angular';
import { BackendSrv } from 'app/core/services/backend_srv';
import { TemplateSrv } from 'app/features/templating/template_srv';

export class GraphiteDatasource {
  basicAuth: string;
  url: string;
  name: string;
  graphiteVersion: any;
  supportsTags: boolean;
  cacheTimeout: any;
  withCredentials: boolean;
  funcDefs: any = null;
  funcDefsPromise: Promise<any> = null;
  _seriesRefLetters: string;

  /** @ngInject */
  constructor(
    instanceSettings: any,
    private $q: IQService,
    private backendSrv: BackendSrv,
    private templateSrv: TemplateSrv
  ) {
    this.basicAuth = instanceSettings.basicAuth;
    this.url = instanceSettings.url;
    this.name = instanceSettings.name;
    this.graphiteVersion = instanceSettings.jsonData.graphiteVersion || '0.9';
    this.supportsTags = supportsTags(this.graphiteVersion);
    this.cacheTimeout = instanceSettings.cacheTimeout;
    this.withCredentials = instanceSettings.withCredentials;
    this.funcDefs = null;
    this.funcDefsPromise = null;

    this._seriesRefLetters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';
  }

  getQueryOptionsInfo() {
    return {
      maxDataPoints: true,
      cacheTimeout: true,
      links: [
        {
          text: 'Help',
          url: 'http://docs.grafana.org/features/datasources/graphite/#using-graphite-in-grafana',
        },
      ],
    };
  }

  query(options: any) {
    const graphOptions = {
      from: this.translateTime(options.rangeRaw.from, false, options.timezone),
      until: this.translateTime(options.rangeRaw.to, true, options.timezone),
      targets: options.targets,
      format: options.format,
      cacheTimeout: options.cacheTimeout || this.cacheTimeout,
      maxDataPoints: options.maxDataPoints,
    };

    const params = this.buildGraphiteParams(graphOptions, options.scopedVars);
    if (params.length === 0) {
      return this.$q.when({ data: [] });
    }

    const httpOptions: any = {
      method: 'POST',
      url: '/render',
      data: params.join('&'),
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
    };

    this.addTracingHeaders(httpOptions, options);

    if (options.panelId) {
      httpOptions.requestId = this.name + '.panelId.' + options.panelId;
    }

    return this.doGraphiteRequest(httpOptions).then(this.convertDataPointsToMs);
  }

  addTracingHeaders(httpOptions: { headers: any }, options: { dashboardId: any; panelId: any }) {
    const proxyMode = !this.url.match(/^http/);
    if (proxyMode) {
      httpOptions.headers['X-Dashboard-Id'] = options.dashboardId;
      httpOptions.headers['X-Panel-Id'] = options.panelId;
    }
  }

  convertDataPointsToMs(result: any) {
    if (!result || !result.data) {
      return [];
    }
    for (let i = 0; i < result.data.length; i++) {
      const series = result.data[i];
      for (let y = 0; y < series.datapoints.length; y++) {
        series.datapoints[y][1] *= 1000;
      }
    }
    return result;
  }

  parseTags(tagString: string) {
    let tags: string[] = [];
    tags = tagString.split(',');
    if (tags.length === 1) {
      tags = tagString.split(' ');
      if (tags[0] === '') {
        tags = [];
      }
    }
    return tags;
  }

  annotationQuery(options: { annotation: { target: string; tags: string }; rangeRaw: any }) {
    // Graphite metric as annotation
    if (options.annotation.target) {
      const target = this.templateSrv.replace(options.annotation.target, {}, 'glob');
      const graphiteQuery = {
        rangeRaw: options.rangeRaw,
        targets: [{ target: target }],
        format: 'json',
        maxDataPoints: 100,
      };

      return this.query(graphiteQuery).then((result: { data: any[] }) => {
        const list = [];

        for (let i = 0; i < result.data.length; i++) {
          const target = result.data[i];

          for (let y = 0; y < target.datapoints.length; y++) {
            const datapoint = target.datapoints[y];
            if (!datapoint[0]) {
              continue;
            }

            list.push({
              annotation: options.annotation,
              time: datapoint[1],
              title: target.target,
            });
          }
        }

        return list;
      });
    } else {
      // Graphite event as annotation
      const tags = this.templateSrv.replace(options.annotation.tags);
      return this.events({ range: options.rangeRaw, tags: tags }).then((results: any) => {
        const list = [];
        for (let i = 0; i < results.data.length; i++) {
          const e = results.data[i];

          let tags = e.tags;
          if (_.isString(e.tags)) {
            tags = this.parseTags(e.tags);
          }

          list.push({
            annotation: options.annotation,
            time: e.when * 1000,
            title: e.what,
            tags: tags,
            text: e.data,
          });
        }

        return list;
      });
    }
  }

  events(options: { range: any; tags: any; timezone?: any }) {
    try {
      let tags = '';
      if (options.tags) {
        tags = '&tags=' + options.tags;
      }
      return this.doGraphiteRequest({
        method: 'GET',
        url:
          '/events/get_data?from=' +
          this.translateTime(options.range.from, false, options.timezone) +
          '&until=' +
          this.translateTime(options.range.to, true, options.timezone) +
          tags,
      });
    } catch (err) {
      return this.$q.reject(err);
    }
  }

  targetContainsTemplate(target: { target: any }) {
    return this.templateSrv.variableExists(target.target);
  }

  translateTime(date: any, roundUp: any, timezone: any) {
    if (_.isString(date)) {
      if (date === 'now') {
        return 'now';
      } else if (date.indexOf('now-') >= 0 && date.indexOf('/') === -1) {
        date = date.substring(3);
        date = date.replace('m', 'min');
        date = date.replace('M', 'mon');
        return date;
      }
      date = dateMath.parse(date, roundUp, timezone);
    }

    // graphite' s from filter is exclusive
    // here we step back one minute in order
    // to guarantee that we get all the data that
    // exists for the specified range
    if (roundUp) {
      if (date.get('s')) {
        date.add(1, 's');
      }
    } else if (roundUp === false) {
      if (date.get('s')) {
        date.subtract(1, 's');
      }
    }

    return date.unix();
  }

  metricFindQuery(query: string, optionalOptions: any) {
    const options: any = optionalOptions || {};
    const interpolatedQuery = this.templateSrv.replace(query);

    // special handling for tag_values(<tag>[,<expression>]*), this is used for template variables
    let matches = interpolatedQuery.match(/^tag_values\(([^,]+)((, *[^,]+)*)\)$/);
    if (matches) {
      const expressions = [];
      const exprRegex = /, *([^,]+)/g;
      let match = exprRegex.exec(matches[2]);
      while (match !== null) {
        expressions.push(match[1]);
        match = exprRegex.exec(matches[2]);
      }
      options.limit = 10000;
      return this.getTagValuesAutoComplete(expressions, matches[1], undefined, options);
    }

    // special handling for tags(<expression>[,<expression>]*), this is used for template variables
    matches = interpolatedQuery.match(/^tags\(([^,]*)((, *[^,]+)*)\)$/);
    if (matches) {
      const expressions = [];
      if (matches[1]) {
        expressions.push(matches[1]);
        const exprRegex = /, *([^,]+)/g;
        let match = exprRegex.exec(matches[2]);
        while (match !== null) {
          expressions.push(match[1]);
          match = exprRegex.exec(matches[2]);
        }
      }
      options.limit = 10000;
      return this.getTagsAutoComplete(expressions, undefined, options);
    }

    const httpOptions: any = {
      method: 'POST',
      url: '/metrics/find',
      params: {},
      data: `query=${interpolatedQuery}`,
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      // for cancellations
      requestId: options.requestId,
    };

    if (options.range) {
      httpOptions.params.from = this.translateTime(options.range.from, false, options.timezone);
      httpOptions.params.until = this.translateTime(options.range.to, true, options.timezone);
    }

    return this.doGraphiteRequest(httpOptions).then((results: any) => {
      return _.map(results.data, metric => {
        return {
          text: metric.text,
          expandable: metric.expandable ? true : false,
        };
      });
    });
  }

  getTags(optionalOptions: any) {
    const options = optionalOptions || {};

    const httpOptions: any = {
      method: 'GET',
      url: '/tags',
      // for cancellations
      requestId: options.requestId,
    };

    if (options.range) {
      httpOptions.params.from = this.translateTime(options.range.from, false, options.timezone);
      httpOptions.params.until = this.translateTime(options.range.to, true, options.timezone);
    }

    return this.doGraphiteRequest(httpOptions).then((results: any) => {
      return _.map(results.data, tag => {
        return {
          text: tag.tag,
          id: tag.id,
        };
      });
    });
  }

  getTagValues(tag: string, optionalOptions: any) {
    const options = optionalOptions || {};

    const httpOptions: any = {
      method: 'GET',
      url: '/tags/' + this.templateSrv.replace(tag),
      // for cancellations
      requestId: options.requestId,
    };

    if (options.range) {
      httpOptions.params.from = this.translateTime(options.range.from, false, options.timezone);
      httpOptions.params.until = this.translateTime(options.range.to, true, options.timezone);
    }

    return this.doGraphiteRequest(httpOptions).then((results: any) => {
      if (results.data && results.data.values) {
        return _.map(results.data.values, value => {
          return {
            text: value.value,
            id: value.id,
          };
        });
      } else {
        return [];
      }
    });
  }

  getTagsAutoComplete(expressions: any[], tagPrefix: any, optionalOptions: any) {
    const options = optionalOptions || {};

    const httpOptions: any = {
      method: 'GET',
      url: '/tags/autoComplete/tags',
      params: {
        expr: _.map(expressions, expression => this.templateSrv.replace((expression || '').trim())),
      },
      // for cancellations
      requestId: options.requestId,
    };

    if (tagPrefix) {
      httpOptions.params.tagPrefix = tagPrefix;
    }
    if (options.limit) {
      httpOptions.params.limit = options.limit;
    }
    if (options.range) {
      httpOptions.params.from = this.translateTime(options.range.from, false, options.timezone);
      httpOptions.params.until = this.translateTime(options.range.to, true, options.timezone);
    }

    return this.doGraphiteRequest(httpOptions).then((results: any) => {
      if (results.data) {
        return _.map(results.data, tag => {
          return { text: tag };
        });
      } else {
        return [];
      }
    });
  }

  getTagValuesAutoComplete(expressions: any[], tag: any, valuePrefix: any, optionalOptions: any) {
    const options = optionalOptions || {};

    const httpOptions: any = {
      method: 'GET',
      url: '/tags/autoComplete/values',
      params: {
        expr: _.map(expressions, expression => this.templateSrv.replace((expression || '').trim())),
        tag: this.templateSrv.replace((tag || '').trim()),
      },
      // for cancellations
      requestId: options.requestId,
    };

    if (valuePrefix) {
      httpOptions.params.valuePrefix = valuePrefix;
    }
    if (options.limit) {
      httpOptions.params.limit = options.limit;
    }
    if (options.range) {
      httpOptions.params.from = this.translateTime(options.range.from, false, options.timezone);
      httpOptions.params.until = this.translateTime(options.range.to, true, options.timezone);
    }

    return this.doGraphiteRequest(httpOptions).then((results: any) => {
      if (results.data) {
        return _.map(results.data, value => {
          return { text: value };
        });
      } else {
        return [];
      }
    });
  }

  getVersion(optionalOptions: any) {
    const options = optionalOptions || {};

    const httpOptions = {
      method: 'GET',
      url: '/version',
      requestId: options.requestId,
    };

    return this.doGraphiteRequest(httpOptions)
      .then((results: any) => {
        if (results.data) {
          const semver = new SemVersion(results.data);
          return semver.isValid() ? results.data : '';
        }
        return '';
      })
      .catch(() => {
        return '';
      });
  }

  createFuncInstance(funcDef: any, options?: any) {
    return gfunc.createFuncInstance(funcDef, options, this.funcDefs);
  }

  getFuncDef(name: string) {
    return gfunc.getFuncDef(name, this.funcDefs);
  }

  waitForFuncDefsLoaded() {
    return this.getFuncDefs();
  }

  getFuncDefs() {
    if (this.funcDefsPromise !== null) {
      return this.funcDefsPromise;
    }

    if (!supportsFunctionIndex(this.graphiteVersion)) {
      this.funcDefs = gfunc.getFuncDefs(this.graphiteVersion);
      this.funcDefsPromise = Promise.resolve(this.funcDefs);
      return this.funcDefsPromise;
    }

    const httpOptions = {
      method: 'GET',
      url: '/functions',
    };

    this.funcDefsPromise = this.doGraphiteRequest(httpOptions)
      .then((results: any) => {
        if (results.status !== 200 || typeof results.data !== 'object') {
          this.funcDefs = gfunc.getFuncDefs(this.graphiteVersion);
        } else {
          this.funcDefs = gfunc.parseFuncDefs(results.data);
        }
        return this.funcDefs;
      })
      .catch((err: any) => {
        console.log('Fetching graphite functions error', err);
        this.funcDefs = gfunc.getFuncDefs(this.graphiteVersion);
        return this.funcDefs;
      });

    return this.funcDefsPromise;
  }

  testDatasource() {
    const query = {
      panelId: 3,
      rangeRaw: { from: 'now-1h', to: 'now' },
      targets: [{ target: 'constantLine(100)' }],
      maxDataPoints: 300,
    };
    return this.query(query).then(() => {
      return { status: 'success', message: 'Data source is working' };
    });
  }

  doGraphiteRequest(options: {
    method?: string;
    url: any;
    requestId?: any;
    withCredentials?: any;
    headers?: any;
    inspect?: any;
  }) {
    if (this.basicAuth || this.withCredentials) {
      options.withCredentials = true;
    }
    if (this.basicAuth) {
      options.headers = options.headers || {};
      options.headers.Authorization = this.basicAuth;
    }

    options.url = this.url + options.url;
    options.inspect = { type: 'graphite' };

    return this.backendSrv.datasourceRequest(options);
  }

  buildGraphiteParams(options: any, scopedVars: ScopedVars) {
    const graphiteOptions = ['from', 'until', 'rawData', 'format', 'maxDataPoints', 'cacheTimeout'];
    const cleanOptions = [],
      targets: any = {};
    let target, targetValue, i;
    const regex = /\#([A-Z])/g;
    const intervalFormatFixRegex = /'(\d+)m'/gi;
    let hasTargets = false;

    options['format'] = 'json';

    function fixIntervalFormat(match: any) {
      return match.replace('m', 'min').replace('M', 'mon');
    }

    for (i = 0; i < options.targets.length; i++) {
      target = options.targets[i];
      if (!target.target) {
        continue;
      }

      if (!target.refId) {
        target.refId = this._seriesRefLetters[i];
      }

      targetValue = this.templateSrv.replace(target.target, scopedVars);
      targetValue = targetValue.replace(intervalFormatFixRegex, fixIntervalFormat);
      targets[target.refId] = targetValue;
    }

    function nestedSeriesRegexReplacer(match: any, g1: string | number) {
      return targets[g1] || match;
    }

    for (i = 0; i < options.targets.length; i++) {
      target = options.targets[i];
      if (!target.target) {
        continue;
      }

      targetValue = targets[target.refId];
      targetValue = targetValue.replace(regex, nestedSeriesRegexReplacer);
      targets[target.refId] = targetValue;

      if (!target.hide) {
        hasTargets = true;
        cleanOptions.push('target=' + encodeURIComponent(targetValue));
      }
    }

    _.each(options, (value, key) => {
      if (_.indexOf(graphiteOptions, key) === -1) {
        return;
      }
      if (value) {
        cleanOptions.push(key + '=' + encodeURIComponent(value));
      }
    });

    if (!hasTargets) {
      return [];
    }

    return cleanOptions;
  }
}

function supportsTags(version: string): boolean {
  return isVersionGtOrEq(version, '1.1');
}

function supportsFunctionIndex(version: string): boolean {
  return isVersionGtOrEq(version, '1.1');
}
