/** @jsx jsx */
import { jsx, Flex } from "theme-ui"
import { Link, useConfig } from "docz"

import * as styles from "gatsby-theme-docz/src/components/Logo/styles"

export const Logo = () => {
  const config = useConfig()
  return (
    <Flex alignItems="center" sx={styles.logo} data-testid="logo">
      <Link
        to="/"
        sx={{
          ...styles.link,
          display: "flex",
          flexDirection: "column",
          lineHeight: "1.2",
        }}
      >
        <span>{config.title}</span>
        <span sx={{ fontSize: "0.7em" }}>{config.description}</span>
      </Link>
    </Flex>
  )
}
