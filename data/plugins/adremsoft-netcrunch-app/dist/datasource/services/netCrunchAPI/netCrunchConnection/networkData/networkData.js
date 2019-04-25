'use strict';System.register(['./networkAtlas'],function(_export){'use strict';var NetCrunchNetworkAtlas;function NetCrunchNetworkData(a,b){function c(l,m,n,o){var p=new a.RemoteDataListStore('ncSrv',1e3,g),q=this;return new Promise(function(r){'function'==typeof n&&p.on('record-changed',function(s){null!=p.data&&0<p.data.length&&s.forEach(n,q)}),'function'==typeof o&&p.on('changed',function(){o()}),p.open(l,m,function(){r()})})}function d(l){h.addNode(l)}function e(l){h.addMap(l)}function f(l){return{accessProfileId:l.AccessProfile,orgId:l.OrgId}}var g=b.serverConnection,h=new NetCrunchNetworkAtlas(g),i={},j={},k=null;return i.promise=new Promise(function(l){return i.resolve=l}),j.promise=new Promise(function(l){return j.resolve=l}),{nodes:function(){return i.promise},networks:function(){return j.promise},atlas:function(){return Promise.all([i.promise,j.promise]).then(function(){return h})},init:function(){var t,u,n=this,_f=f(b.userProfile),o=_f.accessProfileId,p=_f.orgId;return null==k?(t=c('Hosts','Select Id, Name, Address, DeviceType, GlobalDataNode where CanAccessNode(Id, \''+o+':'+p+'\')',d,function(){i.resolve(h.nodes),'function'==typeof n.onNodesChanged&&n.onNodesChanged()}),u=c('Networks','Select NetIntId, DisplayName, HostMapData, IconId, MapType, NetworkData, MapClassTag where (MapClassTag != \'pnet\') && (MapClassTag != \'dependencynet\') && (MapClassTag != \'issuesnet\') && (MapClassTag != \'all\') && (NetIntId != '+2+')',e,function(){j.resolve(h.atlasMaps),'function'==typeof n.onNetworksChanged&&n.onNetworksChanged()}),k=Promise.all([t,u]),k):k}}}/**
 * @license
 * Copyright AdRem Software. All Rights Reserved.
 *
 * Use of this source code is governed by an Apache License, Version 2.0 that can be
 * found in the LICENSE file.
 *//* eslint-disable func-names, object-shorthand, prefer-template */return{setters:[function(_networkAtlas){NetCrunchNetworkAtlas=_networkAtlas.NetCrunchNetworkAtlas}],execute:function(){_export('NetCrunchNetworkData',NetCrunchNetworkData)}}});
//# sourceMappingURL=networkData.js.map
