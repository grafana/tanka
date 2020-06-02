import React from "react"

const Highlight = ({ title, children }) => (
  <div
    css={{
      flexGrow: 1,
      flexBasis: "300px",
      marginTop: "0",
      marginBottom: "1em",
      marginLeft: "2em",
    }}
  >
    <h3 css={{ marginTop: 0, marginBottom: ".5em" }}>{title}</h3>
    {children}
  </div>
)

export const Highlights = ({ elems }) => (
  <div
    css={{
      display: "flex",
      flexDirection: "row",
      width: "100%",
      flexWrap: "wrap",
      marginLeft: "-2em",
      marginBottom: "2em",
    }}
  >
    {Object.keys(elems).map(k => (
      <Highlight key={k} title={k}>
        {elems[k]}
      </Highlight>
    ))}
  </div>
)
