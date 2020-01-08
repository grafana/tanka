/** @jsx jsx */
import { jsx } from "theme-ui"
import React from "react"

import * as headings from "gatsby-theme-docz/src/components/Headings"
import { Layout } from "gatsby-theme-docz/src/components/Layout"
import { Playground } from "gatsby-theme-docz/src/components/Playground"
import { Props } from "gatsby-theme-docz/src/components/Props"
import ThemeStyles from "gatsby-theme-docz/src/theme/styles"

import "typeface-fira-mono"
import "typeface-source-sans-pro"

import { Code, CodeBlock, Pre } from "./codeblock"

const localStyles = {
  backgroundLight: "#2d37471f",
}

// custom "box" (blockquote)
const Box = props => (
  <div
    sx={{
      borderLeft: ".25em solid black",
      borderColor: "text",
      padding: ".25em",
      paddingLeft: "1em",
      background: localStyles.backgroundLight,
      marginBottom: "1rem",
    }}
  >
    {// remove the marginBottom from the last element
    React.Children.map(props.children, (child, i) =>
      i === React.Children.toArray(props.children).length - 1
        ? React.cloneElement(child, {
            style: { marginBottom: 0 },
          })
        : child
    )}
  </div>
)

const Table = props => (
  <div
    sx={{
      overflowX: "auto",
    }}
  >
    <table
      {...props}
      sx={{
        ...ThemeStyles.table,
      }}
    ></table>
  </div>
)

const inlineCode = props => (
  <Code style={{ background: localStyles.backgroundLight }}>
    {props.children}
  </Code>
)

export default {
  ...headings,
  playground: Playground,
  layout: Layout,
  props: Props,
  code: CodeBlock,
  pre: Pre,
  blockquote: Box,
  table: Table,
  inlineCode: inlineCode,
}
