define([
  'lodash'
],
function (_) {
  'use strict';

  function InfluxQueryBuilder(target, database, limit) {
    this.target = target;
    this.database = database;
    this.limit = limit;
  }

  function renderTagCondition (tag, index) {
    var str = "";
    var operator = tag.operator;
    var value = tag.value;
    if (index > 0) {
      str = (tag.condition || 'AND') + ' ';
    }

    if (!operator) {
      if (/^\/.*\/$/.test(tag.value)) {
        operator = '=~';
      } else {
        operator = '=';
      }
    }

    // quote value unless regex or number
    if (operator !== '=~' && operator !== '!~' && isNaN(+value)) {
      value = "'" + value + "'";
    }

    return str + '"' + tag.key + '" ' + operator + ' ' + value;
  }

  var p = InfluxQueryBuilder.prototype;

  p.build = function() {
    return this.target.rawQuery ? this._modifyRawQuery() : this._buildQuery();
  };

  p.buildExploreQuery = function(type, withKey, withMeasurementFilter) {
    var query;
    var measurement;

    if (type === 'TAG_KEYS') {
      query = 'SHOW TAG KEYS';
      measurement = this.target.measurement;
    } else if (type === 'TAG_VALUES') {
      query = 'SHOW TAG VALUES';
      measurement = this.target.measurement;
    } else if (type === 'MEASUREMENTS') {
      query = 'SHOW MEASUREMENTS';
      if (withMeasurementFilter)
      {
        query += ' WITH MEASUREMENT =~ /' + withMeasurementFilter +'/';
      }
    } else if (type === 'FIELDS') {
      if (!this.target.measurement.match('^/.*/')) {
        return 'SHOW FIELD KEYS FROM "' + this.target.measurement + '"';
      } else {
        return 'SHOW FIELD KEYS FROM ' + this.target.measurement;
      }
    } else if (type === 'RETENTION POLICIES') {
      query = 'SHOW RETENTION POLICIES on "' + this.database + '"';
      return query;
    }

    if (measurement) {
      if (!measurement.match('^/.*/') && !measurement.match(/^merge\(.*\)/)) {
        measurement = '"' + measurement+ '"';
      }
      query += ' FROM ' + measurement;
    }

    if (withKey) {
      query += ' WITH KEY = "' + withKey + '"';
    }

    if (this.target.tags && this.target.tags.length > 0) {
      var whereConditions = _.reduce(this.target.tags, function(memo, tag) {
        // do not add a condition for the key we want to explore for
        if (tag.key === withKey) {
          return memo;
        }
        memo.push(renderTagCondition(tag, memo.length));
        return memo;
      }, []);

      if (whereConditions.length > 0) {
        query +=  ' WHERE ' + whereConditions.join(' ');
      }
    }
    if (type === 'MEASUREMENTS')
    {
      query += ' LIMIT ' + (this.limit || 100);
      //Solve issue #2524 by limiting the number of measurements returned
      //LIMIT must be after WITH MEASUREMENT and WHERE clauses
      //This also could be used for TAG KEYS and TAG VALUES, if desired
    }
    return query;
  };

  return InfluxQueryBuilder;
});
