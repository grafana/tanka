/** @jsx jsx */
import { jsx } from "theme-ui"
import React from "react"
import { useDocs } from "docz"

const spacing = "7em"

const Button = ({ href, names, next, alone }) => (
  <a
    href={href}
    sx={{
      flexShrink: 0,
      flexGrow: 1,

      marginLeft: spacing,
      display: "flex",
      flexDirection: "column",
      textAlign: next && !alone ? "right" : "left",
      textDecoration: "none",
      "&:visited": {
        color: "primary",
      },
    }}
  >
    <span sx={{ color: "gray" }}>{next ? "Next" : "Previous"}</span>
    <span
      sx={{
        fontWeight: 700,
      }}
    >
      {names[href]}
    </span>
  </a>
)

export default ({ prev, next }) => {
  const names = useDocs().reduce(
    (map, val) => ({ ...map, [val.route]: val.name }),
    {}
  )

  return (
    <>
      <hr sx={{ marginTop: "4em" }}></hr>
      <div
        sx={{
          display: "flex",
          flexDirection: "row",
          marginLeft: `-${spacing}`,
          flexWrap: "wrap",
          justifyContent: "flex-start",
        }}
      >
        {prev && <Button alone={!next} names={names} href={prev}></Button>}
        {next && <Button next alone={!prev} names={names} href={next}></Button>}
      </div>
    </>
  )
}
