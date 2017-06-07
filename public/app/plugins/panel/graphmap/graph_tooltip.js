/*! grafana - v4.0.2-1481203731 - 2016-12-08
 * Copyright (c) 2016 Torkel Ödegaard; Licensed Apache-2.0 */

define(["jquery","lodash"],function(a){"use strict";function b(b,c,d,e){var f=this,g=d.ctrl,h=g.panel,i=a('<div id="tooltip" class="graph-tooltip">');this.destroy=function(){i.remove()},this.findHoverIndexFromDataPoints=function(a,b,c){for(var d=b.datapoints.pointsize,e=c*d,f=b.datapoints.points.length,g=e;g<f;g+=d)if(!b.lines.steps&&null!=b.datapoints.points[e]&&null==b.datapoints.points[g]||b.datapoints.points[g]>a)return Math.max(g-d,0)/d;return g/d-1},this.findHoverIndexFromData=function(a,b){for(var c=b.data.length,d=0;d<c;d++)if(b.data[d][0]>a)return Math.max(d-1,0);return d-1},this.showTooltip=function(a,b,c,d){"time"===d&&(b='<div class="graph-tooltip-time">'+a+"</div>"+b),i.html(b).place_tt(c.pageX+20,c.pageY)},this.getMultiSeriesPlotHoverInfo=function(a,b){var c,d,e,f,g,i,j,k,l,m=[[],[],[]],n=0;for(d=0;d<a.length;d++)e=a[d],!e.data.length||h.legend.hideEmpty&&e.allIsNull?m[0].push({hidden:!0,value:0}):!e.data.length||h.legend.hideZero&&e.allIsZero?m[0].push({hidden:!0,value:0}):(f=this.findHoverIndexFromData(b.x,e),g=b.x-e.data[f][0],i=e.data[f][0],(!k||g>=0&&(g<k||k<0)||g<0&&g>k)&&(k=g,l=i),e.stack?"individual"===h.tooltip.value_type?c=e.data[f][1]:e.stack?(n+=e.data[f][1],c=n):c=e.data[f][1]:c=e.data[f][1],(e.lines.steps||e.stack)&&(f=this.findHoverIndexFromDataPoints(b.x,e,f)),j=0,e.yaxis&&(j=e.yaxis.n),m[j].push({value:c,hoverIndex:f,color:e.color,label:e.label,time:i,distance:g,index:d}));return m=m[0].concat(m[1],m[2]),m.time=l,m},b.mouseleave(function(){if(h.tooltip.shared){var a=b.data().plot;a&&(i.detach(),a.unhighlight())}c.sharedCrosshair&&g.publishAppEvent("clearCrosshair")}),b.bind("plothover",function(a,j,k){var l,m,n,o,p,q,r,s,t=b.data().plot,u=t.getData(),v=t.getXAxes(),w=v[0].options.mode,x=e();if(c.sharedCrosshair&&g.publishAppEvent("setCrosshair",{pos:j,scope:d}),0!==x.length)if(s=x[0].hasMsResolution?"YYYY-MM-DD HH:mm:ss.SSS":"YYYY-MM-DD HH:mm:ss",h.tooltip.shared){t.unhighlight();var y=f.getMultiSeriesPlotHoverInfo(u,j);for(r="",n=c.formatDate(y.time,s),2===h.tooltip.sort?y.sort(function(a,b){return b.value-a.value}):1===h.tooltip.sort&&y.sort(function(a,b){return a.value-b.value}),p=0;p<y.length;p++)if(o=y[p],!o.hidden){var z="";k&&o.index===k.seriesIndex&&(z="graph-tooltip-list-item--highlight"),q=x[o.index],m=q.formatValue(o.value),r+='<div class="graph-tooltip-list-item '+z+'"><div class="graph-tooltip-series-name">',r+='<i class="fa fa-minus" style="color:'+o.color+';"></i> '+o.label+":</div>",r+='<div class="graph-tooltip-value">'+m+"</div></div>",t.highlight(o.index,o.hoverIndex)}f.showTooltip(n,r,j,w)}else k?(q=x[k.seriesIndex],l='<div class="graph-tooltip-list-item"><div class="graph-tooltip-series-name">',l+='<i class="fa fa-minus" style="color:'+k.series.color+';"></i> '+q.label+":</div>",m=h.stack&&"individual"===h.tooltip.value_type?k.datapoint[1]-k.datapoint[2]:k.datapoint[1],m=q.formatValue(m),n=c.formatDate(k.datapoint[0],s),l+='<div class="graph-tooltip-value">'+m+"</div>",f.showTooltip(n,l,j,w)):i.detach()})}return b});