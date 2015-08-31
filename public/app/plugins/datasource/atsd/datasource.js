define([
    'angular',
    'lodash',
    'kbn',
    'moment',
    './directives',
    './queryCtrl',
],
function (angular, _, kbn) {
    'use strict';

    var module = angular.module('grafana.services');

    module.factory('AtsdDatasource', function($q, backendSrv, templateSrv) {

        function AtsdDatasource(datasource) {
            this.type = 'atsd';
            this.url = datasource.url;
            this.name = datasource.name;
            this.supportMetrics = true;
            this.basicAuth = datasource.basicAuth;
            
            this.cache = 300000;
            
            this.recent = {};
            this.tagRes = {};
        }

        AtsdDatasource.prototype.changeCache = function(cache) {
            this.cache = (cache !== undefined && cache > 0) ? cache * 1000 : this.cache;
            console.log('cache interval: ' + this.cache / 1000 + 's');
        }
        
        AtsdDatasource.prototype.dropCache = function() {
            this.recent = {};
            this.tagRes = {};
            console.log('cache dropped');
        }

        AtsdDatasource.prototype.query = function(options) {
            console.log('actually getting data');
            console.log(options);

            var start = convertToAtsdTime(options.range.from);
            var end = convertToAtsdTime(options.range.to);
            var qs = [];

            _.each(options.targets, function(target) {
                target.disconnect = options.targets[0].disconnect;
                qs.push(convertTargetToQuery(target, options.interval));
            });

            var queries = _.compact(qs);

            if (_.isEmpty(queries)) {
                var d = $q.defer();
                d.resolve({ data: [] });
                return d.promise;
            }

            var groupByTags = {};
            _.each(queries, function(query) {
                _.each(query.tags, function(val, key) {
                    groupByTags[key] = true;
                });
            });

            return this.performTimeSeriesQuery(queries, start, end).then(function(response) {
            
                console.log('status: ' + JSON.stringify(response.status));
                
                if (response.data === undefined) {
                    return { data: [] };
                }
            
                var disconnect = queries[0].disconnect;
                
                var result = _.map(response.data, function(metricData) {
                    return transformMetricData(metricData, disconnect);
                })[0];

                result.sort(function(a, b) {
                    var nameA = a.target.toLowerCase();
                    var nameB = b.target.toLowerCase();
                    
                    if (nameA < nameB) {
                        return -1;
                    } else if (nameA > nameB) {
                        return 1;
                    } else {
                        return 0;
                    }
                });
                
                return { data: result };
            });
        };

        AtsdDatasource.prototype.performTimeSeriesQuery = function(queries, start, end) {
            var jq = [];

            _.each(queries, function(query) {
                console.log('query: ' + JSON.stringify(query));
                
                if (query.entity !== '' && query.metric !== '') {
                    if (query.implicit) {
                        if (query.tagCombos !== undefined) {
                            _.each(query.tagCombos, function(group, ind) {
                                if (group.en) {
                                    var tg = {};
                                    
                                    _.each(group.data, function(value, key){
                                        tg[key] = [value];
                                    });

                                    jq.push(
                                        {
                                            startDate: start,
                                            endDate: end,
                                            limit: 10000,
                                            entity: query.entity,
                                            metric: query.metric,
                                            tags: tg,
                                            aggregate: {
                                                type: query.statistic.toUpperCase(),
                                                interval: {
                                                    count : query.period.count,
                                                    unit: query.period.unit
                                                }
                                            }
                                        }
                                    );
                                }
                            });
                        }
                    } else {
                        var tg = {};
                        
                        _.each(query.tags, function(value, key){
                            tg[key] = [value];
                        });

                        jq.push(
                            {
                                startDate: start,
                                endDate: end,
                                limit: 10000,
                                entity: query.entity,
                                metric: query.metric,
                                tags: tg,
                                aggregate: {
                                    type: query.statistic.toUpperCase(),
                                    interval: {
                                        count : query.period.count,
                                        unit: query.period.unit
                                    }
                                }
                            }
                        );
                    }
                }
            });
            
            if (jq.length == 0) {
                var d = $q.defer();
                d.resolve({ data: undefined });
                return d.promise;
            }
            
            console.log(JSON.stringify({queries: jq}));

            var options = {
                method: 'POST',
                url: this.url + '/api/v1/series',
                data : {
                    queries: jq
                },
                headers: {
                    Authorization: this.basicAuth
                }
            };
            
            var blob = new Blob([JSON.stringify(options)], {type: "text/plain;charset=utf-8"});
            saveAs(blob, "local.log");

            return backendSrv.datasourceRequest(options).then(function(result){
                return result;
            });
        };

        AtsdDatasource.prototype.suggestEntities = function(query) { //DONE
            var so = this;
        
            if (!('entities' in so.recent) ||
                (new Date()).getTime() - so.recent['entities'].time > so.cache) {
                
                so.recent['entities'] = {
                    time: (new Date()).getTime(),
                    value: []
                };
            
                var options = {
                    method: 'GET',
                    url: so.url + '/api/v1/entities',
                    headers: {
                        Authorization: so.basicAuth
                    }
                };

                return backendSrv.datasourceRequest(options).then(function(result) {
                    if (result.status !== 200) {
                        delete so.recent['entities']; 
                        return [];
                    }
                
                    var names = [];
                    _.each(result.data, function(entity){
                        names.push(entity.name);
                    });
                    
                    console.log('entities: ' + JSON.stringify(names));
                    so.recent['entities'].value = names;
                    
                    return names;
                });
            } else {
                var d = $q.defer();
                d.resolve(so.recent['entities'].value);
                return d.promise;
            }
        };

        AtsdDatasource.prototype.suggestMetrics = function(entity, query) { //DONE
            var so = this;
        
            entity = entity !== undefined ? entity : '';
            
            var key = entity !== '' ?
                      'entities/' + entity + '/metrics' :
                      'metrics';
            
            if (!(key in so.recent) ||
                (new Date()).getTime() - so.recent[key].time > so.cache) {
            
                so.recent[key] = {
                    time: (new Date()).getTime(),
                    value: []
                };
                
                var options = {
                    method: 'GET',
                    url: so.url + '/api/v1/' + key,
                    headers: {
                        Authorization: so.basicAuth
                    }
                };

                return backendSrv.datasourceRequest(options).then(function(result) {
                    if (result.status !== 200) {
                        delete so.recent[key]; 
                        return [];
                    }
                
                    var names = [];
                    _.each(result.data, function(metric){
                        names.push(metric.name);
                    });
                    
                    console.log('metrics: ' + JSON.stringify(names));
                    so.recent[key].value = names;
                    
                    return names;
                });
                
            } else {
                var d = $q.defer();
                d.resolve(so.recent[key].value);
                return d.promise;
            }
        };
        
        AtsdDatasource.prototype.suggestNextSegment = function(entity, segments) { //DONE
            segments = segments !== undefined ? segments : [];
            var query = segments.length > 0 ? segments.join('.') + '.' : '';

            return this.suggestMetrics(entity, query).then(function(names){
                var tokens = [];

                tokens = _.map(names, function(name){
                    return name.substr(query.length, name.length).split('.')[0];
                });

                tokens = tokens.filter(function(elem, pos) {
                    return tokens.indexOf(elem) === pos;
                });

                return tokens;
            });

        };
        
        AtsdDatasource.prototype.queryTags = function(entity, metric) {
            var so = this;
        
            entity = entity !== undefined ? entity : '';
            metric = metric !== undefined ? metric : '';
            
            if (entity === '' || metric === '') {
                var d = $q.defer();
                d.resolve({ data: {} });
                return d.promise;
            }
            
            if (!(metric in so.tagRes) ||
                !(entity in so.tagRes[metric]) ||
                (new Date()).getTime() - so.tagRes[metric][entity].time > so.cache) {
                
                if (!(metric in so.tagRes)) {
                    so.tagRes[metric] = {};
                }
                
                so.tagRes[metric][entity] = {
                    time: (new Date()).getTime(),
                    value: {}
                };
                
                var options = {
                    method: 'GET',
                    url: so.url + '/api/v1/metrics/' + metric.replace('%', '%25') + '/entity-and-tags',
                    headers: {
                        Authorization: so.basicAuth
                    }
                };

                return backendSrv.datasourceRequest(options).then(function(result) {
                    if (result.status !== 200) {
                        delete so.tagRes[metric][entity]; 
                        return { data: {} };
                    }
                
                    so.tagRes[metric][entity].value = result;
                    return result;
                });
            } else {
                var d = $q.defer();
                d.resolve(so.tagRes[metric][entity].value);
                return d.promise;
            }
        };
        
        AtsdDatasource.prototype.suggestTags = function(entity, metric, tags_known) { //DONE
            tags_known = tags_known !== undefined ? tags_known : {};
            
            return this.queryTags(entity, metric).then(function(result) {
                console.log('tag query: ' + JSON.stringify(result));
            
                var tags = {};
                _.each(result.data, function(entry) {
                    if (entry.entity === entity) {
                        var matched = true;
                    
                        _.each(entry.tags, function(value, key) {
                            if (key in tags_known && value !== tags_known[key]) {
                                matched = false;
                            }
                        });
                        
                        if (matched) {
                            _.each(entry.tags, function(value, key) {
                                if (!(key in tags_known)) {
                                    if (key in tags) {
                                        tags[key].push(value);
                                    } else {
                                        tags[key] = [value];
                                    }
                                }
                            });
                        }
                    }
                });

                _.each(tags, function(values, key) {
                    tags[key] = values.filter(function(elem, pos) {
                        return values.indexOf(elem) === pos;
                    });
                });

                return tags;
            });
        };

        AtsdDatasource.prototype.suggestTagKeys = function(entity, metric, tags_known) { //DONE
            return this.suggestTags(entity, metric, tags_known).then(function(tags) {
                var keys = _.map(tags, function(values, key) { return key; });
                console.log('tag keys: ' + JSON.stringify(keys));
                return keys;
            });
        };

        AtsdDatasource.prototype.suggestTagValues = function(entity, metric, tags_known, name) { //DONE
            name = name !== undefined ? name : '';
            
            return this.suggestTags(entity, metric, tags_known).then(function(tags) {
                console.log('tag values: ' + JSON.stringify(tags[name]));
                return tags[name];
            });
        };
        
        AtsdDatasource.prototype.getTags = function(entity, metric) { //DONE
            return this.queryTags(entity, metric).then(function(result) {
                var tags = [];
                _.each(result.data, function(entry) {
                    if (entry.entity === entity) {
                        tags.push(entry.tags);
                    }
                });
                console.log('tags: ' + JSON.stringify(tags));
                return tags;
            });
        };

        AtsdDatasource.prototype.testDatasource = function() { //DONE
            var options = {
                method: 'POST',
                url: this.url + '/api/v1/series',
                data : {
                    queries: []
                },
                headers: {
                    Authorization: this.basicAuth
                }
            };
        
            return backendSrv.datasourceRequest(options).then(function(result) {
                return { status: "success", message: "Data source is working", title: "Success" };
            });
        };

        function transformMetricData(metricData, disconnect) {
            var dps;
            var ret = [];

            _.each(metricData, function(dataset) {
                dps = [];
                
                if (disconnect > 0) {
                    if (dataset.data.length > 0) {
                        dps.push([dataset.data[0].v, dataset.data[0].t]);
                        
                        for (var i = 1; i < dataset.data.length; i++) {
                            if (dataset.data[i].t - dataset.data[i-1].t > disconnect * 1000) {
                                dps.push([null, dataset.data[i-1].t + 1]);
                                dps.push([null, dataset.data[i].t - 1]);
                            }
                            
                            dps.push([dataset.data[i].v, dataset.data[i].t]);
                        }
                    }
                } else {
                    _.each(dataset.data, function(data) {
                        dps.push([data.v, data.t]);
                    });
                }

                var name = dataset.entity + ': ' + dataset.metric;

                _.each(dataset.tags, function(value, key) {
                    name += ', ' + key + '=' + value;
                });

                ret.push({ target: name, datapoints: dps });
            });

            return ret;
        }
        
        function parsePeriod(period) {
            var count = '';
            var unit;
            
            for (var i = 0; i < period.length; i++) {
                var c = period.charAt(i);

                if (!isNaN(c)) {
                    count += c;
                } else {
                    unit = c;
                    
                    switch (unit) {
                        case 'y':
                            unit = 'YEAR';
                            break;
                        case 'M':
                            unit = 'MONTH';
                            break;
                        case 'w':
                            unit = 'WEEK';
                            break;
                        case 'd':
                            unit = 'DAY';
                            break;
                        case 'h':
                        case 'H':
                            unit = 'HOUR';
                            break;
                        case 'm':
                            unit = 'MINUTE';
                            break;
                        case 's':
                            unit = 'SECOND';
                            break;
                        default:
                            unit = '';
                    }
                    
                    break;
                }
            }
            
            return { count: parseInt(count), unit: unit };
        };

        function convertTargetToQuery(target, interval) {
            if (!target.metric || !target.entity || target.hide) {
                return null;
            }

            console.log(JSON.stringify(target.period));

            var query = {
                entity: templateSrv.replace(target.entity),
                metric: templateSrv.replace(target.metric),
                
                statistic: target.statistic !== undefined ? templateSrv.replace(target.statistic) : 'detail',
                period: (target.period !== undefined && target.period !== '') ? parsePeriod(target.period) : { count: 1, unit: 'DAY'},
                
                tags: angular.copy(target.tags),
                tagCombos: angular.copy(target.tagCombos),
                implicit: angular.copy(target.implicit),
                
                disconnect: (target.disconnect !== undefined && target.disconnect !== '') ?
                            angular.copy(target.disconnect) :
                            0
            };
            
            console.log(JSON.stringify(query.period));

            return query;
        }

        function convertToAtsdTime(date) { //DONE
            date = date !== 'now' ? date : new Date();
            date = kbn.parseDate(date);
            
            return date.toISOString();
        }

        return AtsdDatasource;
    });
    
});
