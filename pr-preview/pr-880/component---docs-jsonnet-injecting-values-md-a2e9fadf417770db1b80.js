(window.webpackJsonp=window.webpackJsonp||[]).push([[16],{"3fvT":function(e,a,t){"use strict";t.r(a),t.d(a,"_frontmatter",(function(){return r})),t.d(a,"default",(function(){return o}));var n=t("wx14"),s=t("Ff2n"),c=(t("q1tI"),t("7ljp")),m=t("hhGP");t("qKvR");const l=["components"],r={};void 0!==r&&r&&r===Object(r)&&Object.isExtensible(r)&&!Object.prototype.hasOwnProperty.call(r,"__filemeta")&&Object.defineProperty(r,"__filemeta",{configurable:!0,value:{name:"_frontmatter",filename:"docs/jsonnet/injecting-values.md"}});const p={_frontmatter:r},b=m.a;function o(e){let{components:a}=e,t=Object(s.a)(e,l);return Object(c.b)(b,Object(n.a)({},p,t,{components:a,mdxType:"MDXLayout"}),Object(c.b)("h1",{id:"injecting-values"},"Injecting Values"),Object(c.b)("p",null,"Sometimes it might be required to pass externally acquired data into Jsonnet."),Object(c.b)("p",null,"There are three ways of doing so:"),Object(c.b)("ol",null,Object(c.b)("li",{parentName:"ol"},Object(c.b)("a",{parentName:"li",href:"#json-files"},"JSON files")),Object(c.b)("li",{parentName:"ol"},Object(c.b)("a",{parentName:"li",href:"#external-variables"},"External variables")),Object(c.b)("li",{parentName:"ol"},Object(c.b)("a",{parentName:"li",href:"#top-level-arguments"},"Top level arguments"))),Object(c.b)("p",null,"Also check out the ",Object(c.b)("a",{parentName:"p",href:"https://jsonnet.org/ref/language.html#passing-data-to-jsonnet"},"official Jsonnet docs on this\ntopic"),"."),Object(c.b)("h2",{id:"json-files"},"JSON files"),Object(c.b)("p",null,"Jsonnet is a superset of JSON, it treats any JSON as valid Jsonnet. Because many\nsystems can be told to output their data in JSON format, this provides a pretty\ngood interface between those."),Object(c.b)("p",null,"For example, your build tooling like ",Object(c.b)("inlineCode",{parentName:"p"},"make")," could acquire secrets from systems such as\n",Object(c.b)("a",{parentName:"p",href:"https://www.vaultproject.io/"},"Vault"),", etc. and write that into ",Object(c.b)("inlineCode",{parentName:"p"},"secrets.json"),"."),Object(c.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"jsonnet"},Object(c.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk4"},"local"),Object(c.b)("span",{parentName:"span",className:"mtk1"}," secrets "),Object(c.b)("span",{parentName:"span",className:"mtk12"},"="),Object(c.b)("span",{parentName:"span",className:"mtk1"}," "),Object(c.b)("span",{parentName:"span",className:"mtk4"},"import"),Object(c.b)("span",{parentName:"span",className:"mtk1"}," "),Object(c.b)("span",{parentName:"span",className:"mtk6"},'"secrets.json"'),Object(c.b)("span",{parentName:"span",className:"mtk1"},";")),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"})),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"},"{")),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(c.b)("span",{parentName:"span",className:"mtk10"},"foo:"),Object(c.b)("span",{parentName:"span",className:"mtk1"}," secrets.myPassword,")),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"},"}")))),Object(c.b)("blockquote",null,Object(c.b)("p",{parentName:"blockquote"},Object(c.b)("strong",{parentName:"p"},"Note"),": Using ",Object(c.b)("inlineCode",{parentName:"p"},"import")," with JSON treats it as Jsonnet, so make sure to not\nuse it with untrusted code.",Object(c.b)("br",{parentName:"p"}),"\n","A safer, but more verbose, alternative is ",Object(c.b)("inlineCode",{parentName:"p"},"std.parseJson(importstr 'path_to_json.json')"))),Object(c.b)("h2",{id:"external-variables"},"External variables"),Object(c.b)("p",null,"Another way of passing values from the outside are external variables, which are specified like so:"),Object(c.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"bash"},Object(c.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk3 mtki"},"# strings")),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"},"$ tk show "),Object(c.b)("span",{parentName:"span",className:"mtk9"},"."),Object(c.b)("span",{parentName:"span",className:"mtk1"}," --ext-str hello=world")),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"})),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk3 mtki"},"# any Jsonnet snippet")),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"},"$ tk show "),Object(c.b)("span",{parentName:"span",className:"mtk9"},"."),Object(c.b)("span",{parentName:"span",className:"mtk1"}," --ext-code foo=4 --ext-code bar="),Object(c.b)("span",{parentName:"span",className:"mtk5"},"'"),Object(c.b)("span",{parentName:"span",className:"mtk6"},"[ 1, 3 ]"),Object(c.b)("span",{parentName:"span",className:"mtk5"},"'")))),Object(c.b)("p",null,"They can be accessed using ",Object(c.b)("inlineCode",{parentName:"p"},"std.extVar")," and the name given to them on the command line:"),Object(c.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"jsonnet"},Object(c.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"},"{")),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(c.b)("span",{parentName:"span",className:"mtk10"},"foo:"),Object(c.b)("span",{parentName:"span",className:"mtk1"}," std."),Object(c.b)("span",{parentName:"span",className:"mtk9"},"extVar"),Object(c.b)("span",{parentName:"span",className:"mtk1"},"("),Object(c.b)("span",{parentName:"span",className:"mtk6"},"'foo'"),Object(c.b)("span",{parentName:"span",className:"mtk1"},"), "),Object(c.b)("span",{parentName:"span",className:"mtk3 mtki"},"// 4, integer")),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(c.b)("span",{parentName:"span",className:"mtk10"},"bar:"),Object(c.b)("span",{parentName:"span",className:"mtk1"}," std."),Object(c.b)("span",{parentName:"span",className:"mtk9"},"extVar"),Object(c.b)("span",{parentName:"span",className:"mtk1"},"("),Object(c.b)("span",{parentName:"span",className:"mtk6"},"'bar'"),Object(c.b)("span",{parentName:"span",className:"mtk1"},"), "),Object(c.b)("span",{parentName:"span",className:"mtk3 mtki"},"// [ 1, 3 ], array")),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"},"}")))),Object(c.b)("blockquote",null,Object(c.b)("p",{parentName:"blockquote"},Object(c.b)("strong",{parentName:"p"},"Warning"),": External variables are directly accessible in all parts of the\nconfiguration, which can make it difficult to track where they are used and\nwhat effect they have on the final result.\nTry to use ",Object(c.b)("a",{parentName:"p",href:"#top-level-arguments"},"Top Level Arguments")," instead.")),Object(c.b)("h2",{id:"top-level-arguments"},"Top Level Arguments"),Object(c.b)("p",null,"Usually with Tanka, your ",Object(c.b)("inlineCode",{parentName:"p"},"main.jsonnet")," holds an object at the top level (most\nouter type in the generated JSON):"),Object(c.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"jsonnet"},Object(c.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk3 mtki"},"// main.jsonnet")),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"},"{")),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(c.b)("span",{parentName:"span",className:"mtk3 mtki"},"/* your resources */")),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"},"}")))),Object(c.b)("p",null,"Another type of Jsonnet that naturally accepts parameters is the ",Object(c.b)("inlineCode",{parentName:"p"},"function"),".\nWhen the Jsonnet compiler finds a function at the top level, it invokes it and\nallows passing parameter values from the command line:"),Object(c.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"jsonnet"},Object(c.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk3 mtki"},"// Actual output (object) returned by function, which is taking parameters and default values")),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk9"},"function"),Object(c.b)("span",{parentName:"span",className:"mtk1"},"(who, msg="),Object(c.b)("span",{parentName:"span",className:"mtk6"},'"Hello %s!"'),Object(c.b)("span",{parentName:"span",className:"mtk1"},") {")),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(c.b)("span",{parentName:"span",className:"mtk10"},"hello:"),Object(c.b)("span",{parentName:"span",className:"mtk1"}," msg "),Object(c.b)("span",{parentName:"span",className:"mtk12"},"%"),Object(c.b)("span",{parentName:"span",className:"mtk1"}," who")),"\n",Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"},"}")))),Object(c.b)("p",null,"Here, ",Object(c.b)("inlineCode",{parentName:"p"},"who")," needs a value while ",Object(c.b)("inlineCode",{parentName:"p"},"msg")," has a default. This can be invoked like so:"),Object(c.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"bash"},Object(c.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(c.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(c.b)("span",{parentName:"span",className:"mtk1"},"$ tk show "),Object(c.b)("span",{parentName:"span",className:"mtk9"},"."),Object(c.b)("span",{parentName:"span",className:"mtk1"}," --tla-str who=John")))),Object(c.b)("style",{className:"vscode-highlight-styles"},"\n  \n  .material-theme-darker {\nbackground-color: #212121;\ncolor: #EEFFFF;\n}\n\n.material-theme-darker .mtk1 { color: #FFFFFF; }\n.material-theme-darker .mtk2 { color: #212121; }\n.material-theme-darker .mtk3 { color: #545454; }\n.material-theme-darker .mtk4 { color: #F78C6C; }\n.material-theme-darker .mtk5 { color: #89DDFF; }\n.material-theme-darker .mtk6 { color: #C3E88D; }\n.material-theme-darker .mtk7 { color: #FFCB6B; }\n.material-theme-darker .mtk8 { color: #EEFFFF; }\n.material-theme-darker .mtk9 { color: #82AAFF; }\n.material-theme-darker .mtk10 { color: #FF5370; }\n.material-theme-darker .mtk11 { color: #F07178; }\n.material-theme-darker .mtk12 { color: #C792EA; }\n.material-theme-darker .mtk13 { color: #EEFFFF90; }\n.material-theme-darker .mtk14 { color: #65737E; }\n.material-theme-darker .mtk15 { color: #B2CCD6; }\n.material-theme-darker .mtk16 { color: #C17E70; }\n.material-theme-darker .mtki { font-style: italic; }\n.material-theme-darker .mtkb { font-weight: bold; }\n.material-theme-darker .mtku { text-decoration: underline; text-underline-position: under; }\n"))}void 0!==o&&o&&o===Object(o)&&Object.isExtensible(o)&&!Object.prototype.hasOwnProperty.call(o,"__filemeta")&&Object.defineProperty(o,"__filemeta",{configurable:!0,value:{name:"MDXContent",filename:"docs/jsonnet/injecting-values.md"}}),o.isMDXComponent=!0}}]);
//# sourceMappingURL=component---docs-jsonnet-injecting-values-md-a2e9fadf417770db1b80.js.map