/*! grafana - v4.0.2-1481203731 - 2016-12-08
 * Copyright (c) 2016 Torkel Ödegaard; Licensed Apache-2.0 */

System.register(["app/core/utils/kbn"],function(a){function b(){"use strict";return{restrict:"E",scope:!0,templateUrl:"public/app/plugins/panel/graph/axes_editor.html",controller:d}}var c,d;return a("axesEditorComponent",b),{setters:[function(a){c=a}],execute:function(){d=function(){function a(a,b){this.$scope=a,this.$q=b,this.panelCtrl=a.ctrl,this.panel=this.panelCtrl.panel,a.ctrl=this,this.unitFormats=c["default"].getUnitFormats(),this.logScales={linear:1,"log (base 2)":2,"log (base 10)":10,"log (base 32)":32,"log (base 1024)":1024},this.xAxisModes={Time:"time",Series:"series"},this.xAxisStatOptions=[{text:"Avg",value:"avg"},{text:"Min",value:"min"},{text:"Max",value:"min"},{text:"Total",value:"total"},{text:"Count",value:"count"},{text:"Current",value:"current"}],"custom"===this.panel.xaxis.mode&&(this.panel.xaxis.name||(this.panel.xaxis.name="specify field"))}return a.$inject=["$scope","$q"],a.prototype.setUnitFormat=function(a,b){a.format=b.value,this.panelCtrl.render()},a.prototype.render=function(){this.panelCtrl.render()},a.prototype.xAxisOptionChanged=function(){this.panel.xaxis.values&&this.panel.xaxis.values[0]||this.panelCtrl.processor.setPanelDefaultsForNewXAxisMode(),this.panelCtrl.onDataReceived(this.panelCtrl.dataList)},a.prototype.getDataFieldNames=function(a){var b=this.panelCtrl.processor.getDataFieldNames(this.panelCtrl.dataList,a),c=b.map(function(a){return{text:a,value:a}});return this.$q.when(c)},a}(),a("AxesEditorCtrl",d)}}});