import{r as s,b as d,j as e,F as p,a as u,aF as g,aG as o}from"./index.49192afe.js";import{C as m}from"./index.b723ad44.js";function x({err:r}){const[i,c]=s.exports.useState([{id:0,actionName:"",actionUser:"",createAt:0}]);return s.exports.useEffect(()=>{d.get("/eapi/v1/log/top5").then(a=>{const{code:t,msg:n,data:l}=a.data;t!=="0"?r(n):c(l)}).catch(()=>{}).finally(()=>{})},[]),e(p,{children:u(m,{children:[e(g.Title,{heading:5,children:"\u64CD\u4F5C\u8BB0\u5F55"}),e(o,{style:{},size:"small",dataSource:i,render:(a,t)=>e(o.Item,{children:a.actionUser+"->"+a.actionName},t)})]})})}export{x as default};
