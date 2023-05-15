(window.webpackJsonp=window.webpackJsonp||[]).push([[23],{pI5B:function(a,e,t){"use strict";t.r(e),t.d(e,"_frontmatter",(function(){return c})),t.d(e,"default",(function(){return i}));var s=t("wx14"),n=t("zLVn"),m=(t("q1tI"),t("7ljp")),p=t("BIGe"),r=(t("qKvR"),["components"]),c={};void 0!==c&&c&&c===Object(c)&&Object.isExtensible(c)&&!Object.prototype.hasOwnProperty.call(c,"__filemeta")&&Object.defineProperty(c,"__filemeta",{configurable:!0,value:{name:"_frontmatter",filename:"docs/libraries/install-publish.md"}});var b={_frontmatter:c},l=p.a;function i(a){var e=a.components,t=Object(n.a)(a,r);return Object(m.b)(l,Object(s.a)({},b,t,{components:e,mdxType:"MDXLayout"}),Object(m.b)("h1",{id:"installing-and-publishing"},"Installing and publishing"),Object(m.b)("p",null,"The tool for dealing with libraries is\n",Object(m.b)("a",{parentName:"p",href:"https://github.com/jsonnet-bundler/jsonnet-bundler"},Object(m.b)("inlineCode",{parentName:"a"},"jsonnet-bundler")),". It can\ninstall packages from any git source using ",Object(m.b)("inlineCode",{parentName:"p"},"ssh")," and GitHub over ",Object(m.b)("inlineCode",{parentName:"p"},"https"),"."),Object(m.b)("h2",{id:"install-a-library"},"Install a library"),Object(m.b)("p",null,"To install a library from GitHub, use one of the following:"),Object(m.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"bash"},Object(m.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"$ jb install github.com/"),Object(m.b)("span",{parentName:"span",className:"mtk12"},"<"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"user"),Object(m.b)("span",{parentName:"span",className:"mtk12"},">"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"/"),Object(m.b)("span",{parentName:"span",className:"mtk12"},"<"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"repo"),Object(m.b)("span",{parentName:"span",className:"mtk12"},">")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"$ jb install github.com/"),Object(m.b)("span",{parentName:"span",className:"mtk12"},"<"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"user"),Object(m.b)("span",{parentName:"span",className:"mtk12"},">"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"/"),Object(m.b)("span",{parentName:"span",className:"mtk12"},"<"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"repo"),Object(m.b)("span",{parentName:"span",className:"mtk12"},">"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"/"),Object(m.b)("span",{parentName:"span",className:"mtk12"},"<"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"subdir"),Object(m.b)("span",{parentName:"span",className:"mtk12"},">")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"$ jb install github.com/"),Object(m.b)("span",{parentName:"span",className:"mtk12"},"<"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"user"),Object(m.b)("span",{parentName:"span",className:"mtk12"},">"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"/"),Object(m.b)("span",{parentName:"span",className:"mtk12"},"<"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"repo"),Object(m.b)("span",{parentName:"span",className:"mtk12"},">"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"/"),Object(m.b)("span",{parentName:"span",className:"mtk12"},"<"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"subdir"),Object(m.b)("span",{parentName:"span",className:"mtk12"},">"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"@"),Object(m.b)("span",{parentName:"span",className:"mtk12"},"<"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"version"),Object(m.b)("span",{parentName:"span",className:"mtk12"},">")))),Object(m.b)("p",null,"Otherwise, use the ssh syntax:"),Object(m.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"bash"},Object(m.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"$ jb install git+ssh://git@mycode.server:"),Object(m.b)("span",{parentName:"span",className:"mtk12"},"<"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"path-to-repo"),Object(m.b)("span",{parentName:"span",className:"mtk12"},">"),Object(m.b)("span",{parentName:"span",className:"mtk1"},".git")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"$ jb install git+ssh://git@mycode.server:"),Object(m.b)("span",{parentName:"span",className:"mtk12"},"<"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"path-to-repo"),Object(m.b)("span",{parentName:"span",className:"mtk12"},">"),Object(m.b)("span",{parentName:"span",className:"mtk1"},".git/"),Object(m.b)("span",{parentName:"span",className:"mtk12"},"<"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"subdir"),Object(m.b)("span",{parentName:"span",className:"mtk12"},">")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"$ jb install git+ssh://git@mycode.server:"),Object(m.b)("span",{parentName:"span",className:"mtk12"},"<"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"path-to-repo"),Object(m.b)("span",{parentName:"span",className:"mtk12"},">"),Object(m.b)("span",{parentName:"span",className:"mtk1"},".git/"),Object(m.b)("span",{parentName:"span",className:"mtk12"},"<"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"subdir"),Object(m.b)("span",{parentName:"span",className:"mtk12"},">"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"@"),Object(m.b)("span",{parentName:"span",className:"mtk12"},"<"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"version"),Object(m.b)("span",{parentName:"span",className:"mtk12"},">")))),Object(m.b)("blockquote",null,Object(m.b)("p",{parentName:"blockquote"},Object(m.b)("strong",{parentName:"p"},"Note"),": ",Object(m.b)("inlineCode",{parentName:"p"},"version")," may be any git ref, such as commits, tags or branches")),Object(m.b)("h2",{id:"publish-to-github"},"Publish to Git(Hub)"),Object(m.b)("p",null,"Publishing is as easy as committing and pushing to a git remote.\n",Object(m.b)("a",{parentName:"p",href:"https://github.com"},"GitHub")," is recommended, as it is most common and supports\nfaster installing using http archives."),Object(m.b)("style",{className:"vscode-highlight-styles"},"\n  \n  .material-theme-darker {\nbackground-color: #212121;\ncolor: #EEFFFF;\n}\n\n.material-theme-darker .mtk1 { color: #FFFFFF; }\n.material-theme-darker .mtk2 { color: #212121; }\n.material-theme-darker .mtk3 { color: #545454; }\n.material-theme-darker .mtk4 { color: #F78C6C; }\n.material-theme-darker .mtk5 { color: #89DDFF; }\n.material-theme-darker .mtk6 { color: #C3E88D; }\n.material-theme-darker .mtk7 { color: #FFCB6B; }\n.material-theme-darker .mtk8 { color: #EEFFFF; }\n.material-theme-darker .mtk9 { color: #82AAFF; }\n.material-theme-darker .mtk10 { color: #FF5370; }\n.material-theme-darker .mtk11 { color: #F07178; }\n.material-theme-darker .mtk12 { color: #C792EA; }\n.material-theme-darker .mtk13 { color: #EEFFFF90; }\n.material-theme-darker .mtk14 { color: #65737E; }\n.material-theme-darker .mtk15 { color: #B2CCD6; }\n.material-theme-darker .mtk16 { color: #C17E70; }\n.material-theme-darker .mtki { font-style: italic; }\n.material-theme-darker .mtkb { font-weight: bold; }\n.material-theme-darker .mtku { text-decoration: underline; text-underline-position: under; }\n"))}void 0!==i&&i&&i===Object(i)&&Object.isExtensible(i)&&!Object.prototype.hasOwnProperty.call(i,"__filemeta")&&Object.defineProperty(i,"__filemeta",{configurable:!0,value:{name:"MDXContent",filename:"docs/libraries/install-publish.md"}}),i.isMDXComponent=!0}}]);
//# sourceMappingURL=component---docs-libraries-install-publish-md-aed788c89f3c64917380.js.map