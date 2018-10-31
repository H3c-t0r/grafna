import _ from 'lodash';
import angular from 'angular';
import moment from 'moment';

import * as rangeUtil from 'app/core/utils/rangeutil';
import * as customRangeCtrl from '../custom_time_ranges/range_ctrl';

export class TimePickerCtrl {
  static tooltipFormat = 'MMM D, YYYY HH:mm:ss';
  static defaults = {
    time_options: ['5m', '15m', '1h', '6h', '12h', '24h', '2d', '7d', '30d'],
    refresh_intervals: ['5s', '10s', '30s', '1m', '5m', '15m', '30m', '1h', '2h', '1d'],
  };

  dashboard: any;
  panel: any;
  absolute: any;
  timeRaw: any;
  editTimeRaw: any;
  tooltip: string;
  rangeString: string;
  customRangeString: string;
  timeOptions: any;
  customTimeOptions: any;
  refresh: any;
  isUtc: boolean;
  firstDayOfWeek: number;
  isOpen: boolean;
  isAbsolute: boolean;
  isCustom: boolean;
  customRangeIndex: any;
  dayShift: any;

  /** @ngInject */
  constructor(private $scope, private $rootScope, private timeSrv) {
    this.$scope.ctrl = this;

    $rootScope.onAppEvent('shift-time-forward', () => this.move(1), $scope);
    $rootScope.onAppEvent('shift-time-backward', () => this.move(-1), $scope);
    $rootScope.onAppEvent('closeTimepicker', this.openDropdown.bind(this), $scope);

    this.dashboard.on('refresh', this.onRefresh.bind(this), $scope);

    // init options
    this.panel = this.dashboard.timepicker;
    _.defaults(this.panel, TimePickerCtrl.defaults);
    this.firstDayOfWeek = moment.localeData().firstDayOfWeek();

    // init time stuff
    this.onRefresh();
  }

  onRefresh() {
    const time = angular.copy(this.timeSrv.timeRange());
    const timeRaw = angular.copy(time.raw);

    if (!this.dashboard.isTimezoneUtc()) {
      time.from.local();
      time.to.local();
      if (moment.isMoment(timeRaw.from)) {
        timeRaw.from.local();
      }
      if (moment.isMoment(timeRaw.to)) {
        timeRaw.to.local();
      }
      this.isUtc = false;
    } else {
      this.isUtc = true;
    }

    this.isCustom
      ? (this.rangeString = this.customRangeString)
      : (this.rangeString = rangeUtil.describeTimeRange(timeRaw));
    this.absolute = { fromJs: time.from.toDate(), toJs: time.to.toDate() };
    this.tooltip = this.dashboard.formatDate(time.from) + ' <br>to<br>';
    this.tooltip += this.dashboard.formatDate(time.to);
    this.timeRaw = timeRaw;
    this.isAbsolute = moment.isMoment(this.timeRaw.to);
  }

  zoom(factor) {
    this.$rootScope.appEvent('zoom-out', 2);
  }

  move(direction) {
    if (!this.isCustom) {
      const range = this.timeSrv.timeRange();

      const timespan = (range.to.valueOf() - range.from.valueOf()) / 2;
      let to, from;
      if (direction === -1) {
        to = range.to.valueOf() - timespan;
        from = range.from.valueOf() - timespan;
      } else if (direction === 1) {
        to = range.to.valueOf() + timespan;
        from = range.from.valueOf() + timespan;
        if (to > Date.now() && range.to < Date.now()) {
          to = Date.now();
          from = range.from.valueOf();
        }
      } else {
        to = range.to.valueOf();
        from = range.from.valueOf();
      }
      this.timeSrv.setTime({ from: moment.utc(from), to: moment.utc(to) });
    } else {
      const functionResult = customRangeCtrl.customMove(
        direction,
        this.customRangeIndex,
        this.customTimeOptions,
        this.dayShift
      );
      this.dayShift = functionResult.dayShift;
      this.customTimeOptionMoved(this.customTimeOptions[functionResult.index], this.dayShift);
    }
  }

  openDropdown() {
    if (this.isOpen) {
      this.closeDropdown();
      return;
    }

    this.onRefresh();
    this.editTimeRaw = this.timeRaw;
    this.timeOptions = rangeUtil.getRelativeTimesList(this.panel, this.rangeString);
    this.customTimeOptions = this.dashboard.ranges;
    this.refresh = {
      value: this.dashboard.refresh,
      options: _.map(this.panel.refresh_intervals, (interval: any) => {
        return { text: interval, value: interval };
      }),
    };

    this.refresh.options.unshift({ text: 'off' });
    this.isOpen = true;
    this.$rootScope.appEvent('timepickerOpen');
  }

  closeDropdown() {
    this.isOpen = false;
    this.$rootScope.appEvent('timepickerClosed');
  }

  applyCustom() {
    if (this.refresh.value !== this.dashboard.refresh) {
      this.timeSrv.setAutoRefresh(this.refresh.value);
    }

    this.timeSrv.setTime(this.editTimeRaw);
    this.closeDropdown();
  }

  absoluteFromChanged() {
    this.editTimeRaw.from = this.getAbsoluteMomentForTimezone(this.absolute.fromJs);
    this.isCustom = false;
  }

  absoluteToChanged() {
    this.editTimeRaw.to = this.getAbsoluteMomentForTimezone(this.absolute.toJs);
    this.isCustom = false;
  }

  getAbsoluteMomentForTimezone(jsDate) {
    return this.dashboard.isTimezoneUtc() ? moment(jsDate).utc() : moment(jsDate);
  }

  setRelativeFilter(timespan) {
    const range = { from: timespan.from, to: timespan.to };

    if (this.panel.nowDelay && range.to === 'now') {
      range.to = 'now-' + this.panel.nowDelay;
    }

    this.timeSrv.setTime(range);
    this.closeDropdown();
    this.isCustom = false;
  }

  // dayShift a shift nejsou nejstastnejsí nazvy :/ pak jeste dayShift a this.dayShift
  // TODO Picked neni uplne korektni slovo treba seleted nebo tak pouyit
  customTimeOptionPicked(range) {
    const time = angular.copy(this.timeSrv.timeRange());
    this.editTimeRaw.from = this.dashboard.formatDate(time.from);
    this.dayShift = customRangeCtrl.customTimeRangePicked('shiftByDay', range, 0, this.editTimeRaw).dayShift;
    this.applyCustomRange(range);
  }

  customTimeOptionMoved(range, dayShift) {
    this.dayShift = dayShift;
    customRangeCtrl.customTimeRangePicked('shift', range, dayShift, this.editTimeRaw);
    this.applyCustomRange(range);
  }

  applyCustomRange(range) {
    this.editTimeRaw.from = this.getAbsoluteMomentForTimezone(range.absoluteFrom);
    this.editTimeRaw.to = this.getAbsoluteMomentForTimezone(range.absoluteTo);
    this.applyCustom();
    this.customRangeString = range.name + ', ' + rangeUtil.describeTimeRange(this.editTimeRaw).substring(0, 12);
    this.isCustom = true;
    this.customRangeIndex = this.customTimeOptions.indexOf(range);
  }
}

export function settingsDirective() {
  return {
    restrict: 'E',
    templateUrl: 'public/app/features/dashboard/timepicker/settings.html',
    controller: TimePickerCtrl,
    bindToController: true,
    controllerAs: 'ctrl',
    scope: {
      dashboard: '=',
    },
  };
}

export function timePickerDirective() {
  return {
    restrict: 'E',
    templateUrl: 'public/app/features/dashboard/timepicker/timepicker.html',
    controller: TimePickerCtrl,
    bindToController: true,
    controllerAs: 'ctrl',
    scope: {
      dashboard: '=',
    },
  };
}

angular.module('grafana.directives').directive('gfTimePickerSettings', settingsDirective);
angular.module('grafana.directives').directive('gfTimePicker', timePickerDirective);

import { inputDateDirective } from './input_date';
angular.module('grafana.directives').directive('inputDatetime', inputDateDirective);
