import{r as s,j as e,F as n,ay as a}from"./index.49192afe.js";function c({handleClose:t,errorMessage:r}){return s.exports.useEffect(()=>{const o=setTimeout(()=>{t()},3e3);return()=>clearTimeout(o)},[t]),e(n,{children:e(a,{type:"error",showIcon:!0,closable:!0,onClose:t,content:r,style:{position:"fixed",top:"20px",left:"50%",width:"320px",height:"60px",transform:"translateX(-50%)",zIndex:9999}})})}export{c as M};
