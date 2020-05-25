import { ChevronDown, ChevronUp } from "gatsby-theme-docz/src/components/Icons"
import React, { useState } from "react"

import ThemeStyles from "gatsby-theme-docz/src/theme/styles"
/** @jsx jsx */
import { jsx } from "theme-ui"

export const Code = props => (
  <code
    sx={{
      ...props.style,
      fontFamily: "Fira Mono, monospace",
      fontSize: "1rem",
    }}
  >
    {props.children}
  </code>
)

// Smart codeblock: shows only first 25 lines, if longer an expand button
export const CodeBlock = props => {
  const lines = React.Children.toArray(props.children).reduce((n, c) => {
    if (c?.props?.className === "vscode-highlight-line") {
      return n + 1
    }
    return n
  }, 0)

  return (
    <Code>
      {lines > 20 ? <LongCode>{props.children}</LongCode> : props.children}
    </Code>
  )
}

export const Pre = props => (
  <pre
    {...props}
    sx={{
      ...ThemeStyles.pre,
      fontSize: "1rem",
      lineHeight: 1.4,
      overflowX: "auto",
    }}
  ></pre>
)

// Expandable codeblock
const LongCode = props => {
  const [toggled, setToggled] = useState(false)

  return (
    <>
      {toggled
        ? props.children
        : React.Children.map(props.children, (child, i) => {
            if (i < 20 * 2) return child
          })}
      <Expand toggled={toggled} onClick={() => setToggled(!toggled)}></Expand>
    </>
  )
}

// ExpandButton
const Expand = props => (
  <button
    onClick={props.onClick}
    sx={{
      background: "inherit",
      border: "none",
      color: "inherit",
      fontFamily: "inherit",
      fontSize: "inherit",
      textDecoration: "underline",
      cursor: "pointer",
      ":hover": {
        textDecoration: "none",
      },
      display: "flex",
      padding: 0,
      width: "100%",
      justifyContent: "center",
    }}
  >
    <div
      sx={{
        display: "flex",
        alignItems: "center",
      }}
    >
      {props.toggled ? <ChevronUp></ChevronUp> : <ChevronDown></ChevronDown>}
      Show {props.toggled ? "less" : "more"}
    </div>
  </button>
)
