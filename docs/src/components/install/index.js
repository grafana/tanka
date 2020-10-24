import React, { useState } from "react"

import Osx from "./tk/osx.mdx"
import Arch from "./tk/arch.mdx"
import Go from "./tk/go.mdx"
import Binary from "./tk/binary.mdx"

import JbOsx from "./jb/osx.mdx"
import JbArch from "./jb/arch.mdx"
import JbGo from "./jb/go.mdx"
import JbBinary from "./jb/binary.mdx"

export const Tanka = {
  macOS: <Osx />,
  ArchLinux: <Arch />,
  Binary: <Binary />,
  Go: <Go />,
}

export const Jb = {
  macOS: <JbOsx />,
  ArchLinux: <JbArch />,
  Binary: <JbBinary />,
  Go: <JbGo />,
}

export const PlatformInstall = ({ elems, def }) => {
  const [current, setCurrent] = useState(def)

  return (
    <div css={{ display: "flex", flexDirection: "column" }}>
      <div
        css={{
          display: "flex",
          marginBottom: "1em",
          marginLeft: "-.5em",
          flexWrap: "wrap",
        }}
      >
        {Object.keys(elems).map(e => (
          <button
            key={e}
            css={{
              background: "none",
              color: "inherit",
              fontSize: "1em",
              fontFamily: "inherit",
              padding: ".4em .8em .4em .8em",
              marginLeft: ".5em",
              marginBottom: ".5em",
              border: `1px solid ${(e === current && "#0B5FFF") || "#CED4DE"}`,
              borderRadius: "3px",
              outline: "none",
              ":hover": {
                border: "1px solid #0B5FFF",
                cursor: "pointer",
              },
            }}
            onClick={() => {
              setCurrent(e)
            }}
          >
            {e}
          </button>
        ))}
      </div>
      <div
        css={{
          marginTop: "-1em",
          border: "1px solid #CED4DE",
          borderRadius: "3px",
          padding: ".8em",
        }}
      >
        {elems[current]}
      </div>
    </div>
  )
}
