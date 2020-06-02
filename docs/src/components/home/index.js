import React from "react"
import { Link } from "docz"

import Snip from "./snippet.mdx"

const elemCss = {
  display: "flex",
  flexGrow: 1,
  flexDirection: "column",
  flexBasis: `calc(50% - 2em)`,
  justifyContent: "center",
  marginBottom: "1em",
  marginLeft: "2em",
}

export const Catcher = () => (
  <div
    css={{
      display: "flex",
      flexWrap: "wrap",
      marginLeft: "-2em",
      marginTop: "1em",
    }}
  >
    <div css={elemCss}>
      <h1
        css={{
          marginTop: 0,
          marginBottom: 0,
          fontSize: "3em",
          lineHeight: "normal",
        }}
      >
        Define. Reuse. Override.
      </h1>
      <p>
        Grafana Tanka is the robust configuration utility for your{" "}
        <a href="https://kubernetes.io">Kubernetes</a> cluster, powered by the
        unique <a href="https://jsonnet.org">Jsonnet</a> language
      </p>
      <div css={{ display: "flex", marginLeft: "-1em" }}>
        <Button to="/install">Install</Button>
        <Button to="/tutorial/overview">Tutorial</Button>
      </div>
    </div>
    <div
      css={{
        ...elemCss,
        flexBasis: "calc(50% - 2em)",
        overflowX: "hidden",
        pre: { marginTop: 0, marginBottom: 0 },
      }}
    >
      <Snip></Snip>
      <small>Kubernetes Deployment. That's all it takes.</small>
    </div>
  </div>
)

const Button = ({ to, children }) => (
  <Link
    css={{
      marginLeft: "1em",
      textDecoration: "none",
      color: "white",
      background: "#0B5FFF",
      padding: ".5em",
      justifyContent: "center",
      display: "flex",
      flexGrow: 1,
      borderRadius: "5px",
      border: "2px solid #0B5FFF",
      ":hover": {
        background: "white",
        color: "#0B5FFF",
      },
    }}
    to={to}
  >
    {children}
  </Link>
)
