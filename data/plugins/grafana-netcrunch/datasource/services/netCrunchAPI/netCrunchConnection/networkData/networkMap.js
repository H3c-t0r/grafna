'use strict';System.register([],function(_export,_context){'use strict';var _createClass,PRIVATE_PROPERTIES,NetCrunchNetworkMap;function _classCallCheck(instance,Constructor){if(!(instance instanceof Constructor))throw new TypeError('Cannot call a class as a function')}return{setters:[],execute:function(){_createClass=function(){function defineProperties(target,props){for(var i=0;i<props.length;i++){var descriptor=props[i];descriptor.enumerable=descriptor.enumerable||!1,descriptor.configurable=!0,'value'in descriptor&&(descriptor.writable=!0),Object.defineProperty(target,descriptor.key,descriptor)}}return function(Constructor,protoProps,staticProps){return protoProps&&defineProperties(Constructor.prototype,protoProps),staticProps&&defineProperties(Constructor,staticProps),Constructor}}();PRIVATE_PROPERTIES={local:Symbol('local'),values:Symbol('values')};_export('NetCrunchNetworkMap',NetCrunchNetworkMap=function(){function NetCrunchNetworkMap(a){/* eslint-disable no-param-reassign */function b(c,d){if(null!=d.HostMapData)for(var e=1,f=d.HostMapData[0];e<=f;e+=1){var g=d.HostMapData[e];(0===g[0]||5===g[0])&&c.nodesId.push(parseInt(g[1],10))}}_classCallCheck(this,NetCrunchNetworkMap),null==a?(this[PRIVATE_PROPERTIES.local]={},this[PRIVATE_PROPERTIES.values]={}):(this[PRIVATE_PROPERTIES.local]=a.local,this[PRIVATE_PROPERTIES.values]=a.getValues()),this[PRIVATE_PROPERTIES.local].nodesId=[],function(d,e){if(d.netId=e.NetIntId||'',d.parentId=null==e.NetworkData?'':parseInt(e.NetworkData[0],10),isNaN(d.parentId)&&(d.parentId=''),d.isFolder='dynfolder'===e.MapClassTag||null!=e.NetworkData&&Array.isArray(e.NetworkData[1]),d.isFolder){var f=null==e.NetworkData?[]:e.NetworkData[1];Array.isArray(f)&&(d.maps=f.map(function(g){return parseInt(g,10)})),'fnet'===e.MapClassTag&&b(d,e)}else b(d,e)}/* eslint-enable no-param-reassign */(this[PRIVATE_PROPERTIES.local],this[PRIVATE_PROPERTIES.values]),this[PRIVATE_PROPERTIES.local].children=[]}return _createClass(NetCrunchNetworkMap,[{key:'addChild',value:function addChild(a){var b=this.children.every(function(c){return c.netId!==a.netId});!0===b&&this.children.push(a)}},{key:'netId',get:function get(){return this[PRIVATE_PROPERTIES.local].netId}},{key:'parentId',get:function get(){return this[PRIVATE_PROPERTIES.local].parentId}},{key:'nodesId',get:function get(){return this[PRIVATE_PROPERTIES.local].nodesId}},{key:'allNodesId',get:function get(){function a(c,d){return c.forEach(function(e){return d.add(e)}),d}var b=new Set;return this.isFolder?('fnet'===this[PRIVATE_PROPERTIES.values].MapClassTag&&a(this.nodesId,b),this.children.forEach(function(c){a(c.allNodesId,b)})):a(this.nodesId,b),Array.from(b)}},{key:'isFolder',get:function get(){return this[PRIVATE_PROPERTIES.local].isFolder}},{key:'children',get:function get(){return this[PRIVATE_PROPERTIES.local].children}},{key:'displayName',get:function get(){return this[PRIVATE_PROPERTIES.values].DisplayName||''}},{key:'allChildren',get:function get(){function a(c,d){return c.displayName.localeCompare(d.displayName)}function b(c,d,e){var f=[];return c.children.sort(a).forEach(function(g){if(f.push({map:g,innerLevel:d,parentIndex:e}),c.isFolder&&2>=d){var h=isNaN(e)?f.length-1:e+f.length;f=f.concat(b(g,d+1,h))}}),f}return b(this,1,'root')}}]),NetCrunchNetworkMap}());_export('NetCrunchNetworkMap',NetCrunchNetworkMap)}}});
//# sourceMappingURL=networkMap.js.map
