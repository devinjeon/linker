(this.webpackJsonpapp=this.webpackJsonpapp||[]).push([[0],{26:function(e,t,n){},27:function(e,t,n){},45:function(e,t,n){"use strict";n.r(t);var a=n(0),c=n(1),s=n.n(c),i=n(14),o=n.n(i),r=(n(26),n(15)),l=n(16),u=n(17),h=n(20),d=n(18),j=n(19),g=n.p+"static/media/logo.6ce24c58.svg",p=(n(27),n(47)),b=n(4),f=n.n(b);function O(e){return e.isLoggedIn?Object(a.jsx)("a",{href:"".concat(f.a.defaults.baseURL,"/auth/signout"),children:"\ub85c\uadf8\uc544\uc6c3"}):Object(a.jsx)("a",{href:"".concat(f.a.defaults.baseURL,"/auth/signin"),children:"\ub85c\uadf8\uc778"})}f.a.defaults.baseURL="https://api.linker.hyojun.me";var m=function(e){Object(h.a)(n,e);var t=Object(d.a)(n);function n(){var e;Object(l.a)(this,n);for(var a=arguments.length,c=new Array(a),s=0;s<a;s++)c[s]=arguments[s];return(e=t.call.apply(t,[this].concat(c))).state={id:"",url:""},e.handleChange=function(t){e.setState(Object(r.a)({},t.target.name,t.target.value))},e.handleSubmit=function(t){t.preventDefault(),e.props.onCreate(e.state).then((function(t){200===e.response.status&&e.setState({id:"",url:""})})).catch((function(t){console.log(e.error)}))},e}return Object(u.a)(n,[{key:"render",value:function(){return Object(a.jsxs)("form",{onSubmit:this.handleSubmit,children:[Object(a.jsx)("input",{placeholder:"ID",value:this.state.id,onChange:this.handleChange,name:"id"}),Object(a.jsx)("input",{placeholder:"URL",value:this.state.url,onChange:this.handleChange,name:"url"}),Object(a.jsx)("input",{type:"submit",value:"\ub4f1\ub85d"})]})}}]),n}(c.Component),v=function(){var e=Object(p.a)(["is_logged_in"]),t=Object(j.a)(e,1)[0];return Object(a.jsx)("div",{className:"App",children:Object(a.jsxs)("header",{className:"App-header",children:[Object(a.jsx)("img",{src:g,className:"App-logo",alt:"logo"}),Object(a.jsx)(m,{onCreate:function(e){return console.log(e),f.a.post("/links/".concat(e.id),{url:e.url},{withCredentials:!0})}}),Object(a.jsx)(O,{isLoggedIn:"true"===t.is_logged_in})]})})},x=function(e){e&&e instanceof Function&&n.e(3).then(n.bind(null,48)).then((function(t){var n=t.getCLS,a=t.getFID,c=t.getFCP,s=t.getLCP,i=t.getTTFB;n(e),a(e),c(e),s(e),i(e)}))};o.a.render(Object(a.jsx)(s.a.StrictMode,{children:Object(a.jsx)(v,{})}),document.getElementById("root")),x(console.log)}},[[45,1,2]]]);
//# sourceMappingURL=main.04dfa803.chunk.js.map