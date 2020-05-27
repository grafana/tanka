/** @jsx jsx */
import { jsx } from "theme-ui"

import * as styles from "gatsby-theme-docz/src/components/NavSearch/styles"
import { Search } from "gatsby-theme-docz/src/components/Icons"

export const Algolia = () => (
  <div
    sx={{
      ...styles.wrapper,
      ".algolia-autocomplete": {
        ...styles.input,
        span: {
          fontSize: "16px",
          fontFamily: "inherit",
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
