import * as styles from "gatsby-theme-docz/src/components/Sidebar/styles"

import { Box, jsx } from "theme-ui"
import { Link, useCurrentDoc, useMenus } from "docz"
/** @jsx jsx */
/** @jsxFrag React.Fragment */
import React, { useEffect, useRef, useState } from "react"

import { Global } from "@emotion/core"
import Logo from "../../../../img/logo.svg"
import { NavGroup } from "gatsby-theme-docz/src/components/NavGroup"
import { NavLink } from "gatsby-theme-docz/src/components/NavLink"
import { NavSearch } from "gatsby-theme-docz/src/components/NavSearch"

import { Algolia } from "./algolia"

export const Sidebar = React.forwardRef((props, ref) => {
  const [query, setQuery] = useState("")
  const menus = useMenus({ query })
  const currentDoc = useCurrentDoc()
  const currentDocRef = useRef()
  const handleChange = ev => {
    setQuery(ev.target.value)
  }
  useEffect(() => {
    if (ref.current && currentDocRef.current) {
      // disabling, so logo stays visible
      // ref.current.scrollTo(0, currentDocRef.current.offsetTop)
    }
  }, [ref])
  return (
    <>
      <Box onClick={props.onClick} sx={styles.overlay(props)}>
        {props.open && <Global styles={styles.global} />}
      </Box>
      <Box ref={ref} sx={styles.wrapper(props)} data-testid="sidebar">
        <Link to="/">
          <img
            src={Logo}
            css={{ width: "100%", marginBottom: "2em", marginTop: "1em" }}
            alt=""
          ></img>
        </Link>

        <Algolia />

        {menus &&
          menus
            .filter(e => {
              if (e.filepath) {
                return e.filepath.startsWith("docs")
              }
              return true
            })
            .map(menu => {
              if (!menu.route)
                return <NavGroup key={menu.id} item={menu} sidebarRef={ref} />
              if (menu.route === currentDoc.route) {
                return (
                  <NavLink key={menu.id} item={menu} ref={currentDocRef}>
                    {menu.name}
                  </NavLink>
                )
              }
              return (
                <NavLink key={menu.id} item={menu}>
                  {menu.name}
                </NavLink>
              )
            })}
      </Box>
    </>
  )
})
