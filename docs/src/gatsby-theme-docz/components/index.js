/** @jsx jsx */
import { jsx } from "theme-ui";
import React from "react";
import { Global } from "@emotion/core";

import * as headings from "gatsby-theme-docz/src/components/Headings";
import { Layout } from "gatsby-theme-docz/src/components/Layout";
import { Playground } from "gatsby-theme-docz/src/components/Playground";
import { Props } from "gatsby-theme-docz/src/components/Props";
import ThemeStyles from "gatsby-theme-docz/src/theme/styles";

import "typeface-fira-mono";

// codeblock style
const Code = props => (
  <code
    sx={{
      fontFamily: "Fira Mono, monospace",
      fontSize: "1em"
    }}
  >
    {props.children}
  </code>
);

const Pre = props => (
  <pre
    {...props}
    sx={{
      ...ThemeStyles.pre,
      fontSize: "1rem",
      lineHeight: 1.4,
      overflowX: "auto"
    }}
  ></pre>
);

// custom "box" (blockquote)
const Box = props => (
  <div
    sx={{
      borderLeft: ".25em solid black",
      borderColor: "text",
      padding: ".25em",
      paddingLeft: "1em",
      background: "#2d37471f"
    }}
  >
    {// remove the marginBottom from the last element
    React.Children.map(props.children, (child, i) =>
      i === React.Children.toArray(props.children).length - 1
        ? React.cloneElement(child, {
            style: { marginBottom: 0 }
          })
        : child
    )}
  </div>
);

const Table = props => (
  <div
    sx={{
      overflowX: "auto"
    }}
  >
    <table
      {...props}
      sx={{
        ...ThemeStyles.table
      }}
    ></table>
  </div>
);

export default {
  ...headings,
  playground: Playground,
  layout: Layout,
  props: Props,
  code: Code,
  pre: Pre,
  blockquote: Box,
  table: Table
};
