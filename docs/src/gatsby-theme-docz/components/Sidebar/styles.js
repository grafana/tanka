import * as styles from "gatsby-theme-docz/src/components/Sidebar/styles"

export const global = styles.global
export const overlay = styles.overlay

export const wrapper = ({ open }) => ({
  ...styles.wrapper(open),
  overflow: "visible",
})
