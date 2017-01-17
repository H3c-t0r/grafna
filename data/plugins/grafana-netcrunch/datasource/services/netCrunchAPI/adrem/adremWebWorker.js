'use strict';System.register([],function(_export,_context){'use strict';var _typeof,_createClass,AdremWebWorker;function _classCallCheck(instance,Constructor){if(!(instance instanceof Constructor))throw new TypeError('Cannot call a class as a function')}return{setters:[],execute:function(){_typeof='function'==typeof Symbol&&'symbol'==typeof Symbol.iterator?function(obj){return typeof obj}:function(obj){return obj&&'function'==typeof Symbol&&obj.constructor===Symbol&&obj!==Symbol.prototype?'symbol':typeof obj};_createClass=function(){function defineProperties(target,props){for(var i=0;i<props.length;i++){var descriptor=props[i];descriptor.enumerable=descriptor.enumerable||!1,descriptor.configurable=!0,'value'in descriptor&&(descriptor.writable=!0),Object.defineProperty(target,descriptor.key,descriptor)}}return function(Constructor,protoProps,staticProps){return protoProps&&defineProperties(Constructor.prototype,protoProps),staticProps&&defineProperties(Constructor,staticProps),Constructor}}();_export('AdremWebWorker',AdremWebWorker=function(){function AdremWebWorker(a){function b(){for(var e=new Date().getTime();c.has(e);)e+=1;return e}_classCallCheck(this,AdremWebWorker);var c=new Map,d=new Worker(a);d.onmessage=function(e){var f=e.data.taskId;if(c.has(f)){var g=c.get(f);c.delete(f),g(e.data.result)}},this.executeTask=function(e){var f=b(),g=e;return g.taskId=f,new Promise(function(h){c.set(f,h),d.postMessage(g)})}}return _createClass(AdremWebWorker,[{key:'addTask',value:function addTask(a){this[a.name]=function(){var _this=this;for(var _len=arguments.length,b=Array(_len),_key=0;_key<_len;_key++)b[_key]=arguments[_key];var c={funcName:a.name,args:b,async:a.async};if(!0===a.async){var _ret=function(){var d=_this;return{v:new Promise(function(e,f){d.executeTask(c).then(function(g){'resolve'===g.type&&e(g.result),'reject'===g.type&&f(g.error)})})}}();if('object'==('undefined'==typeof _ret?'undefined':_typeof(_ret)))return _ret.v}return this.executeTask(c)}}}],[{key:'webWorkerBuilder',value:function webWorkerBuilder(){function a(){function e(){function g(k,l){j.postMessage({taskId:k,result:l})}function h(k,l,m){g(k,j[l].apply(j,m))}function i(k,l,m){j[l].apply(j,m).then(function(n){return g(k,{type:'resolve',result:n})}).catch(function(n){return g(k,{type:'reject',error:n})})}var j=this;return function(l){var m=l.data;!0===m.async?i(m.taskId,m.funcName,m.args):h(m.taskId,m.funcName,m.args)}}var f=function(){return'this.onmessage = '+e.name+'().bind(this);\n\n'}();return f+=e.toString()+'\n',f+=c.reduce(function(g,h){return g+'\n'+h},''),new Blob([f],{type:'application/javascript'})}function b(){return URL.createObjectURL(a())}var c=[],d=[];return{addFunctionCode:function addFunctionCode(f){var g=1<arguments.length&&void 0!==arguments[1]&&arguments[1],h=2<arguments.length&&void 0!==arguments[2]&&arguments[2];return'function'==typeof f&&(c.push(f.toString()),!0===g&&null!=f.name&&''!==f.name&&d.push({name:f.name,async:h}),!0)},getWebWorker:function getWebWorker(){var f=new AdremWebWorker(b());return d.forEach(function(g){f.addTask(g)}),f}}}}]),AdremWebWorker}());_export('AdremWebWorker',AdremWebWorker)}}});
//# sourceMappingURL=adremWebWorker.js.map
