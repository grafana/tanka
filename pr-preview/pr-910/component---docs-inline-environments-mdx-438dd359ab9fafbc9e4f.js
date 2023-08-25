(window.webpackJsonp=window.webpackJsonp||[]).push([[13],{d0WO:function(e,a,n){"use strict";n.r(a),n.d(a,"_frontmatter",(function(){return b})),n.d(a,"default",(function(){return r}));var s=n("wx14"),t=n("Ff2n"),m=(n("q1tI"),n("7ljp")),p=n("hhGP");n("qKvR");const c=["components"],b={};void 0!==b&&b&&b===Object(b)&&Object.isExtensible(b)&&!Object.prototype.hasOwnProperty.call(b,"__filemeta")&&Object.defineProperty(b,"__filemeta",{configurable:!0,value:{name:"_frontmatter",filename:"docs/inline-environments.mdx"}});const l={_frontmatter:b},N=p.a;function r(e){let{components:a}=e,n=Object(t.a)(e,c);return Object(m.b)(N,Object(s.a)({},l,n,{components:a,mdxType:"MDXLayout"}),Object(m.b)("h1",{id:"inline-environments"},"Inline environments"),Object(m.b)("p",null,"Inline environments is the practice of defining the environment's config inline\nfor evaluation at runtime as opposed to configuring it statically in\n",Object(m.b)("inlineCode",{parentName:"p"},"spec.json"),"."),Object(m.b)("p",null,"The general take away is:"),Object(m.b)("ul",null,Object(m.b)("li",{parentName:"ul"},Object(m.b)("inlineCode",{parentName:"li"},"spec.json")," will no longer be used"),Object(m.b)("li",{parentName:"ul"},Object(m.b)("inlineCode",{parentName:"li"},"main.jsonnet")," is expected to render a ",Object(m.b)("inlineCode",{parentName:"li"},"tanka.dev/Environment")," object"),Object(m.b)("li",{parentName:"ul"},"this object is expected to hold Kubernetes objects at ",Object(m.b)("inlineCode",{parentName:"li"},".data"))),Object(m.b)("h2",{id:"converting-to-an-inline-environment"},"Converting to an inline environment"),Object(m.b)("p",null,"Converting a traditional ",Object(m.b)("inlineCode",{parentName:"p"},"spec.json")," environment into an inline environment is quite\nstraight forward. Based on the example from ",Object(m.b)("a",{parentName:"p",href:"tutorial/jsonnet"},"Using Jsonnet"),":"),Object(m.b)("p",null,"The directory structure:"),Object(m.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"sh"},Object(m.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"├── environments")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"│   └── default "),Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"# default environment")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"│       ├── main.jsonnet "),Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"# main file")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"│       └── spec.json "),Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"# environment's config")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"├── jsonnetfile.json")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"├── lib "),Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"# libraries")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"└── vendor "),Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"# external libraries")))),Object(m.b)("p",null,"The original files look like this:"),Object(m.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"jsonnet"},Object(m.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"// main.jsonnet")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"{")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"some_deployment:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," {"),Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"/* ... */"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"some_service:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," {"),Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"/* ... */"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"}")))),Object(m.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"json"},Object(m.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"// spec.json")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk5"},"{")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk12"},"apiVersion"),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk5"},":"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk6"},"tanka.dev/v1alpha1"),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk5"},",")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk12"},"kind"),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk5"},":"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk6"},"Environment"),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk5"},",")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk12"},"metadata"),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk5"},":"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk5"},"{")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk7"},"name"),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk5"},":"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk6"},"default"),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"')),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk5"},"},")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk12"},"spec"),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk5"},":"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk5"},"{")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk7"},"apiServer"),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk5"},":"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk6"},"https://127.0.0.1:6443"),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk5"},",")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk7"},"namespace"),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk5"},":"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"'),Object(m.b)("span",{parentName:"span",className:"mtk6"},"monitoring"),Object(m.b)("span",{parentName:"span",className:"mtk5"},'"')),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk5"},"}")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk5"},"}")))),Object(m.b)("p",null,"Converting is as simple as bringing in the ",Object(m.b)("inlineCode",{parentName:"p"},"spec.json")," into ",Object(m.b)("inlineCode",{parentName:"p"},"main.jsonnet")," and\nmoving the original ",Object(m.b)("inlineCode",{parentName:"p"},"main.jsonnet")," scope into the ",Object(m.b)("inlineCode",{parentName:"p"},"data:")," element."),Object(m.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"jsonnet"},Object(m.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"// main.jsonnet")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"{")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"apiVersion:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'tanka.dev/v1alpha1'"),Object(m.b)("span",{parentName:"span",className:"mtk1"},",")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"kind:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'Environment'"),Object(m.b)("span",{parentName:"span",className:"mtk1"},",")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"metadata:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," {")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"name:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'default'"),Object(m.b)("span",{parentName:"span",className:"mtk1"},",")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"spec:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," {")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"apiServer:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'https://127.0.0.1:6443'"),Object(m.b)("span",{parentName:"span",className:"mtk1"},",")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"namespace:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'monitoring'"),Object(m.b)("span",{parentName:"span",className:"mtk1"},",")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"data:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," { "),Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"// original main.jsonnet data")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"some_deployment:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," {"),Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"/* ... */"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"some_service:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," {"),Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"/* ... */"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"}")))),Object(m.b)("h2",{id:"use-case-variable-apiserver"},"Use case: variable apiServer"),Object(m.b)("p",null,"Even though the ",Object(m.b)("inlineCode",{parentName:"p"},"apiServer")," directive is originally meant to prevent that the\nmanifests don't get accidentally applied to the wrong Kubernetes cluster, there\nis a valid use case for making the ",Object(m.b)("inlineCode",{parentName:"p"},"apiServer")," variable: Local test clusters."),Object(m.b)("p",null,"Instead of modifying ",Object(m.b)("inlineCode",{parentName:"p"},"spec.json")," each time, with inline environments it is\npossible to leverage powerful jsonnet concepts, for example with top level\narguments:"),Object(m.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"jsonnet"},Object(m.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"// environments/minikube-test-setup/main.jsonnet")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk9"},"function"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"(apiServer) {")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"apiVersion:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'tanka.dev/v1alpha1'"),Object(m.b)("span",{parentName:"span",className:"mtk1"},",")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"kind:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'Environment'"),Object(m.b)("span",{parentName:"span",className:"mtk1"},",")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"metadata:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," {")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"name:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'minikube-test-setup'"),Object(m.b)("span",{parentName:"span",className:"mtk1"},",")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"spec:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," {")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"apiServer:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," apiServer,")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"namespace:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'monitoring'"),Object(m.b)("span",{parentName:"span",className:"mtk1"},",")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"data:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," { "),Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"/* ... */"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"}")))),Object(m.b)("p",null,"Applying this to a local Kubernetes cluster can be done like this:"),Object(m.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"bash"},Object(m.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"$ tk apply --tla-str apiServer=https://127.0.0.1:4758 environments/minikube-test-setup")))),Object(m.b)("p",null,"Similarly this can be used to configure any part of the Environment object, like\n",Object(m.b)("inlineCode",{parentName:"p"},"namespace:"),", ",Object(m.b)("inlineCode",{parentName:"p"},"metadata.labels"),", ..."),Object(m.b)("h2",{id:"use-case-consistent-inline-environments"},"Use case: consistent inline environments"),Object(m.b)("p",null,"It is possible to define multiple inline environments in a single jsonnet. This\nenables an operator to generate consistent Tanka environments for multiple\nKubernetes clusters."),Object(m.b)("p",null,"We can define a Tanka environment once and then repeat that for a set of\nclusters as shown in this example:"),Object(m.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"jsonnet"},Object(m.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"// environments/monitoring-stack/main.jsonnet")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"{")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk9"},"environment"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"(cluster)"),Object(m.b)("span",{parentName:"span",className:"mtk12"},"::"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," {")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"apiVersion:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'tanka.dev/v1alpha1'"),Object(m.b)("span",{parentName:"span",className:"mtk1"},",")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"kind:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'Environment'"),Object(m.b)("span",{parentName:"span",className:"mtk1"},",")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"metadata:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," {")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"      "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"name:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'environment/%s'"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk12"},"%"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," cluster.name,")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"spec:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," {")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"      "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"apiServer:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," cluster.apiServer,")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"      "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"namespace:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'monitoring'"),Object(m.b)("span",{parentName:"span",className:"mtk1"},",")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"data:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," { "),Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"/* ... */"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"})),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk7"},"clusters::"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," [")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    { "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"name:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'us-central1'"),Object(m.b)("span",{parentName:"span",className:"mtk1"},", "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"apiServer:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'https://127.0.0.1:6433'"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    { "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"name:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'europe-west2'"),Object(m.b)("span",{parentName:"span",className:"mtk1"},", "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"apiServer:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk6"},"'https://127.0.0.2:6433'"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  ],")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"})),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  "),Object(m.b)("span",{parentName:"span",className:"mtk10"},"envs:"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," {")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    [cluster.name]"),Object(m.b)("span",{parentName:"span",className:"mtk12"},":"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk4"},"$"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"."),Object(m.b)("span",{parentName:"span",className:"mtk9"},"environment"),Object(m.b)("span",{parentName:"span",className:"mtk1"},"(cluster)")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"    "),Object(m.b)("span",{parentName:"span",className:"mtk5 mtki"},"for"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," cluster "),Object(m.b)("span",{parentName:"span",className:"mtk5 mtki"},"in"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," "),Object(m.b)("span",{parentName:"span",className:"mtk4"},"$"),Object(m.b)("span",{parentName:"span",className:"mtk1"},".clusters")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  },")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"}")))),Object(m.b)("p",null,"In the workflow you now have to use ",Object(m.b)("inlineCode",{parentName:"p"},"--name")," to select the environment you want\nto deploy:"),Object(m.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"bash"},Object(m.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"$ tk apply --name environment/us-central1 environments/monitoring-stack/main.jsonnet")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"$ tk diff --name environment/europe-west2 environments/monitoring-stack/main.jsonnet")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"})),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk3 mtki"},"# Partial matches also work (if they match a single environment)")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"$ tk apply --name us-central1 environments/monitoring-stack/main.jsonnet")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"$ tk diff --name west2 environments/monitoring-stack/main.jsonnet")))),Object(m.b)("p",null,"For export, it is possible to use the same ",Object(m.b)("inlineCode",{parentName:"p"},"--name")," selector or you can do a\nrecursive export while using the ",Object(m.b)("inlineCode",{parentName:"p"},"--format")," option:"),Object(m.b)("pre",{className:"material-theme-darker vscode-highlight","data-language":"bash"},Object(m.b)("code",{parentName:"pre",className:"vscode-highlight-code"},Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"$ tk "),Object(m.b)("span",{parentName:"span",className:"mtk12"},"export"),Object(m.b)("span",{parentName:"span",className:"mtk1"}," outputDir/ environments/monitoring-stack/main.jsonnet --recursive \\")),"\n",Object(m.b)("span",{parentName:"code",className:"vscode-highlight-line"},Object(m.b)("span",{parentName:"span",className:"mtk1"},"  --format "),Object(m.b)("span",{parentName:"span",className:"mtk5"},"'"),Object(m.b)("span",{parentName:"span",className:"mtk6"},"{{env.metadata.name}}/{{.metadata.namespace}}/{{.kind}}-{{.metadata.name}}"),Object(m.b)("span",{parentName:"span",className:"mtk5"},"'")))),Object(m.b)("h2",{id:"caveats"},"Caveats"),Object(m.b)("h3",{id:"import-tk"},Object(m.b)("inlineCode",{parentName:"h3"},'import "tk"')),Object(m.b)("p",null,"Inline environments cannot use ",Object(m.b)("a",{parentName:"p",href:"config#jsonnet-access"},Object(m.b)("inlineCode",{parentName:"a"},'import "tk"'))," anymore as\nthis information was populated before jsonnet evaluation by the existence of\n",Object(m.b)("inlineCode",{parentName:"p"},"spec.json"),"."),Object(m.b)("h3",{id:"tk-env"},Object(m.b)("inlineCode",{parentName:"h3"},"tk env")),Object(m.b)("p",null,"The different ",Object(m.b)("inlineCode",{parentName:"p"},"tk env")," subcommands are heavily based on the ",Object(m.b)("inlineCode",{parentName:"p"},"spec.json"),"\napproach. ",Object(m.b)("inlineCode",{parentName:"p"},"tk env list")," will continue to work as expected, ",Object(m.b)("inlineCode",{parentName:"p"},"tk env\n(add|remove|set)")," will only work for ",Object(m.b)("inlineCode",{parentName:"p"},"spec.json")," based environments."),Object(m.b)("style",{className:"vscode-highlight-styles"},"\n  \n  .material-theme-darker {\nbackground-color: #212121;\ncolor: #EEFFFF;\n}\n\n.material-theme-darker .mtk1 { color: #FFFFFF; }\n.material-theme-darker .mtk2 { color: #212121; }\n.material-theme-darker .mtk3 { color: #545454; }\n.material-theme-darker .mtk4 { color: #F78C6C; }\n.material-theme-darker .mtk5 { color: #89DDFF; }\n.material-theme-darker .mtk6 { color: #C3E88D; }\n.material-theme-darker .mtk7 { color: #FFCB6B; }\n.material-theme-darker .mtk8 { color: #EEFFFF; }\n.material-theme-darker .mtk9 { color: #82AAFF; }\n.material-theme-darker .mtk10 { color: #FF5370; }\n.material-theme-darker .mtk11 { color: #F07178; }\n.material-theme-darker .mtk12 { color: #C792EA; }\n.material-theme-darker .mtk13 { color: #EEFFFF90; }\n.material-theme-darker .mtk14 { color: #65737E; }\n.material-theme-darker .mtk15 { color: #B2CCD6; }\n.material-theme-darker .mtk16 { color: #C17E70; }\n.material-theme-darker .mtki { font-style: italic; }\n.material-theme-darker .mtkb { font-weight: bold; }\n.material-theme-darker .mtku { text-decoration: underline; text-underline-position: under; }\n"))}void 0!==r&&r&&r===Object(r)&&Object.isExtensible(r)&&!Object.prototype.hasOwnProperty.call(r,"__filemeta")&&Object.defineProperty(r,"__filemeta",{configurable:!0,value:{name:"MDXContent",filename:"docs/inline-environments.mdx"}}),r.isMDXComponent=!0}}]);
//# sourceMappingURL=component---docs-inline-environments-mdx-438dd359ab9fafbc9e4f.js.map