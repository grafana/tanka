/** @jsx jsx */
import { jsx } from "theme-ui"
import { useEffect } from "react"

import * as styles from "gatsby-theme-docz/src/components/NavSearch/styles"
import { Search } from "gatsby-theme-docz/src/components/Icons"

import config from "../../../../algolia.json"

export const Algolia = () => {
  useEffect(() => {
    if (
      typeof window === "undefined" ||
      typeof window.docsearch === "undefined"
    ) {
      console.log("DocSearch unavailable")
      return
    }

    window.docsearch(config)
  })
  return (
    <div
      sx={{
        ...styles.wrapper,
        ".algolia-autocomplete": {
          ...styles.input,
          ".ds-dropdown-menu": {
            // TODO: not working
            // fontSize: "16px",
          },
        },
        ".algolia-docsearch-suggestion--highlight": {
          color: "primary",
        },
        input: styles.input,
      }}
    >
      <Search size={20} sx={styles.icon} />
      <input id="algolia-docsearch" placeholder="Type to search..." />
    </div>
  )
}
