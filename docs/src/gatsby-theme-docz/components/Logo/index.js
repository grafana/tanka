import * as styles from "gatsby-theme-docz/src/components/Logo/styles"

/** @jsx jsx */
import { Flex, jsx } from "theme-ui"
import { Link, useConfig } from "docz"
import logo from '../../../../img/logo.svg'

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
                <div sx={{ display: 'inline-block' }}>
                    <img src={logo} alt="That's my logo" sx={{ maxHeight: "1.2em", marginRight: '0.2em', display: 'inline-block', float: 'left' }} />
                    <span sx={{ display: 'inline-block' }}>Grafana Tanka</span>
                    <span sx={{ fontSize: "0.7em", display: 'block' }}>{config.description}</span>
                </div>
            </Link>
        </Flex>
    )
}
