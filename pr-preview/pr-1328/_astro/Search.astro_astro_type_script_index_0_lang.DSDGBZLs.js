const w="modulepreload",y=function(m){return"/pr-preview/pr-1328/"+m},g={},S=function(u,i,c){let p=Promise.resolve();if(i&&i.length>0){document.getElementsByTagName("link");const o=document.querySelector("meta[property=csp-nonce]"),t=o?.nonce||o?.getAttribute("nonce");p=Promise.allSettled(i.map(n=>{if(n=y(n),n in g)return;g[n]=!0;const d=n.endsWith(".css"),f=d?'[rel="stylesheet"]':"";if(document.querySelector(`link[href="${n}"]${f}`))return;const r=document.createElement("link");if(r.rel=d?"stylesheet":w,d||(r.as="script"),r.crossOrigin="",r.href=n,t&&r.setAttribute("nonce",t),document.head.appendChild(r),d)return new Promise((e,s)=>{r.addEventListener("load",e),r.addEventListener("error",()=>s(new Error(`Unable to preload CSS for ${n}`)))})}))}function l(o){const t=new Event("vite:preloadError",{cancelable:!0});if(t.payload=o,window.dispatchEvent(t),!t.defaultPrevented)throw o}return p.then(o=>{for(const t of o||[])t.status==="rejected"&&l(t.reason);return u().catch(l)})},E={ranking:{pageLength:.1,termFrequency:.1,termSaturation:2,termSimilarity:9}};class v extends HTMLElement{constructor(){super();const u=this.querySelector("button[data-open-modal]"),i=this.querySelector("button[data-close-modal]"),c=this.querySelector("dialog"),p=this.querySelector(".dialog-frame"),l=e=>{("href"in(e.target||{})||document.body.contains(e.target)&&!p.contains(e.target))&&t()},o=e=>{c.showModal(),document.body.toggleAttribute("data-search-modal-open",!0),this.querySelector("input")?.focus(),e?.stopPropagation(),window.addEventListener("click",l)},t=()=>c.close();u.addEventListener("click",o),u.disabled=!1,i.addEventListener("click",t),c.addEventListener("close",()=>{document.body.toggleAttribute("data-search-modal-open",!1),window.removeEventListener("click",l)}),window.addEventListener("keydown",e=>{(e.metaKey===!0||e.ctrlKey===!0)&&e.key==="k"&&(c.open?t():o(),e.preventDefault())});let n={};try{n=JSON.parse(this.dataset.translations||"{}")}catch{}const r=this.dataset.stripTrailingSlash!==void 0?e=>e.replace(/(.)\/(#.*)?$/,"$1$2"):e=>e;window.addEventListener("DOMContentLoaded",()=>{(window.requestIdleCallback||(s=>setTimeout(s,1)))(async()=>{const{PagefindUI:s}=await S(async()=>{const{PagefindUI:a}=await import("./ui-core.CYBi0lMN.js");return{PagefindUI:a}},[]);new s({...E,element:"#starlight__search",baseUrl:"/pr-preview/pr-1328/",bundlePath:"/pr-preview/pr-1328/".replace(/\/$/,"")+"/pagefind/",showImages:!1,translations:n,showSubResults:!0,processResult:a=>{a.url=r(a.url),a.sub_results=a.sub_results.map(h=>(h.url=r(h.url),h))}})})})}}customElements.define("site-search",v);export{S as _};
