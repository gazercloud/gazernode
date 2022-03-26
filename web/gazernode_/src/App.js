import React, {useState} from 'react';
import PropTypes from 'prop-types';
import AppBar from '@material-ui/core/AppBar';
import CssBaseline from '@material-ui/core/CssBaseline';
import Divider from '@material-ui/core/Divider';
import Drawer from '@material-ui/core/Drawer';
import Hidden from '@material-ui/core/Hidden';
import IconButton from '@material-ui/core/IconButton';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import MenuIcon from '@material-ui/icons/Menu';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import {makeStyles, MuiThemeProvider, useTheme} from '@material-ui/core/styles';
import BlurOnIcon from '@material-ui/icons/BlurOn';
import InfoOutlinedIcon from '@material-ui/icons/InfoOutlined';
import PageUnits from "./Pages/PageUnits";
import PageAbout from "./Pages/PageAbout";
import Grid from "@material-ui/core/Grid";
import PageSensorAdd from "./Pages/PageSensorAdd";
import SignIn from "./Pages/SignIn";
import PageUnit from "./Pages/PageUnit";
import PageAccount from "./Pages/PageAccount";
import {createMuiTheme} from "@material-ui/core";
import PageDataItem from "./Pages/PageDataItem";
import PageUnitConfig from "./Pages/PageUnitConfig";
import PageAddUnit from "./Pages/PageAddUnit";
import PagePublicChannels from "./Pages/PagePublicChannels";
import PagePublicChannel from "./Pages/PagePublicChannel";
import { SnackbarProvider, useSnackbar } from 'notistack';
import PageRemoteAccess from "./Pages/PageRemoteAccess";
import CloudUploadOutlinedIcon from '@material-ui/icons/CloudUploadOutlined';
import CloudOutlinedIcon from '@material-ui/icons/CloudOutlined';
import PeopleIcon from '@material-ui/icons/People';
import PersonIcon from '@material-ui/icons/Person';
import PageUsers from "./Pages/PageUsers";
import PageUser from "./Pages/PageUser";


const drawerWidth = 240;


function getCookie(name) {
    let matches = document.cookie.match(new RegExp(
        "(?:^|; )" + name.replace(/([\.$?*|{}\(\)\[\]\\\/\+^])/g, '\\$1') + "=([^;]*)"
    ));
    return matches ? decodeURIComponent(matches[1]) : undefined;
}


const useStyles = makeStyles((theme) => ({
    root: {
        display: 'flex',
        color: '#D9D9D9',
        backgroundColor: '#121212'
    },
    drawer: {
        [theme.breakpoints.up('sm')]: {
            width: drawerWidth,
            flexShrink: 0
        },
    },
    appBar: {
        [theme.breakpoints.up('sm')]: {
            width: `calc(100% - ${drawerWidth}px)`,
            marginLeft: drawerWidth,
            backgroundColor: "#272727",
            color: "#CCC"
        },
    },
    menuButton: {
        marginRight: theme.spacing(2),
        [theme.breakpoints.up('sm')]: {
            display: 'none'
        },
    },
    // necessary for content to be below app bar
    toolbar: theme.mixins.toolbar,
    drawerPaper: {
        width: drawerWidth,
        color: "#EEEEEE",
        borderRight: '1px solid #850',
        backgroundColor: "#1E1E1E"
    },
    content: {
        flexGrow: 1,
        padding: theme.spacing(3),
    },
}));

function getHashVariable(variable) {
    const query = window.location.hash.substring(1);
    const vars = query.split('&');
    for (let i = 0; i < vars.length; i++) {
        const pair = vars[i].split('=');
        if (decodeURIComponent(pair[0]) === variable) {
            return decodeURIComponent(pair[1]);
        }
    }
    console.log('Query variable %s not found', variable);
}
/*
function getForm(variable) {
    const query = window.location.search.substring(1);
    const vars = query.split('&');
    for (let i = 0; i < vars.length; i++) {
        const pair = vars[i].split('=');
        if (decodeURIComponent(pair[0]) === variable) {
            return decodeURIComponent(pair[1]);
        }
    }
    console.log('Query variable %s not found', variable);
}
*/

const gotoLink = (link) => {
    window.location = link
}

function useForceUpdate(){
    const [value, setValue] = React.useState(0); // integer state
    return () => setValue(value => value + 1); // update the state to force render
}

function getWindow() {
    return window
}

function ResponsiveDrawer(props) {
    const {window} = props;
    const classes = useStyles();
    //const theme = useTheme();
    const [mobileOpen, setMobileOpen] = React.useState(false);
    const [redrawBool, setRedrawBool] = React.useState(false);
    const [title, setTitle] = React.useState(false);

    const forceUpdate = useForceUpdate();

    const navigate = (link) => {
        gotoLink(link)
        setMobileOpen(false)
        setRedrawBool(!redrawBool)
    }

    const updateTitle = (t) => {
        t = "Gazer Node - " + t
        if (title !== t) {
            setTitle(t)
            document.title = t
        }
    }

    const handleDrawerToggle = () => {
        setMobileOpen(!mobileOpen);
        forceUpdate()
    };


    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        setFirstRendering(false)
    }

    getWindow().onhashchange = () => {
        setRedrawBool(!redrawBool)
    }

    const getClientWidth = () => {
        let w = getWindow().innerWidth - 300
        return w
    }

    const drawer = (
        <div>
            <Grid container className={classes.toolbar} direction="row" alignItems="center" alignContent="center">

                <Grid container direction="row" alignItems="flex-start" alignContent="flex-start"
                      style={{margin: "15px"}}>
                    <Grid item><a href="/"><img src="/mainicon32.png" style={{width: "48px"}}/></a></Grid>
                    <Grid item>
                        <Typography style={{marginLeft: "8px", marginBottom: "0px", fontSize: "14pt"}}>
                            GazerNode
                        </Typography>
                        <Typography style={{marginLeft: "8px", fontSize: "10pt", color: "#BBB"}}>
                            monitoring system
                        </Typography>
                    </Grid>
                </Grid>

            </Grid>
            <Divider/>
            <List>
                <ListItem button key="units" component="a" onClick={() => {
                    navigate("#form=units")
                }}>
                    <ListItemIcon><BlurOnIcon style={{color: "#2278B5"}}/></ListItemIcon>
                    <ListItemText primary={"Units"}/>
                </ListItem>
            </List>
            <Divider/>
            <List>
                <ListItem button key="public_channels" component="a" onClick={() => {
                    navigate("#form=public_channels")
                }}>
                    <ListItemIcon><CloudUploadOutlinedIcon style={{color: "#2278B5"}}/></ListItemIcon>
                    <ListItemText primary={"Public Channels"}/>
                </ListItem>
            </List>
            <Divider/>
            <List>
                <ListItem button key="remote_access" component="a" onClick={() => {
                    navigate("#form=remote_access")
                }}>
                    <ListItemIcon><CloudOutlinedIcon style={{color: "#2278B5"}}/></ListItemIcon>
                    <ListItemText primary={"Remote Access"}/>
                </ListItem>
            </List>
            <Divider/>
            <List>
                <ListItem button key="users" component="a" onClick={() => {
                    navigate("#form=users")
                }}>
                    <ListItemIcon><PeopleIcon style={{color: "#2278B5"}}/></ListItemIcon>
                    <ListItemText primary={"Users"}/>
                </ListItem>
            </List>
            <Divider/>
            <List>
                <ListItem button key="account" component="a" onClick={() => {
                    navigate("#form=account")
                }}>
                    <ListItemIcon><PersonIcon style={{color: "#2278B5"}}/></ListItemIcon>
                    <ListItemText primary={"Account"}/>
                </ListItem>
            </List>
            <Divider/>
            <List>
                <ListItem button key="about" component="a" onClick={() => {
                    navigate("#form=about")
                }}>
                    <ListItemIcon><InfoOutlinedIcon style={{color: "#2278B5"}}/></ListItemIcon>
                    <ListItemText primary={"About"}/>
                </ListItem>
            </List>
        </div>
    );

    const renderForm = () => {
        const form = getHashVariable("form")

        if (form === "units") {
            return (
                <PageUnits
                    onAddSensor={() => {
                        navigate("#form=sensor_add")
                    }}
                    OnNavigate={(addr) => navigate(addr)}
                    OnTitleUpdate={(title) => updateTitle(title)}
                />
            )
        }
        if (form === "unit") {
            return (
                <PageUnit
                    OnNavigate={(addr) => navigate(addr)}
                    UnitId={getHashVariable("unitId")}
                    OnTitleUpdate={(title) => updateTitle(title)}
                />
            )
        }
        if (form === "public_channels") {
            return (
                <PagePublicChannels
                    onAddSensor={() => {
                        navigate("#form=sensor_add")
                    }}
                    OnNavigate={(addr) => navigate(addr)}
                    OnTitleUpdate={(title) => updateTitle(title)}
                />
            )
        }
        if (form === "public_channel") {
            return (
                <PagePublicChannel
                    OnNavigate={(addr) => navigate(addr)}
                    ChannelId={getHashVariable("channelId")}
                    OnTitleUpdate={(title) => updateTitle(title)}
                />
            )
        }
        if (form === "users") {
            return (
                <PageUsers
                    OnNavigate={(addr) => navigate(addr)}
                    OnTitleUpdate={(title) => updateTitle(title)}
                />
            )
        }
        if (form === "user") {
            return (
                <PageUser
                    OnNavigate={(addr) => navigate(addr)}
                    UserName={getHashVariable("user")}
                    OnTitleUpdate={(title) => updateTitle(title)}
                />
            )
        }
        if (form === "remote_access") {
            return (
                <PageRemoteAccess
                    OnNavigate={(addr) => navigate(addr)}
                    OnTitleUpdate={(title) => updateTitle(title)}
                />
            )
        }
        if (form === "unit_config") {
            return (
                <PageUnitConfig
                    OnNavigate={(addr) => navigate(addr)}
                    UnitId={getHashVariable("unitId")}
                    UnitType={getHashVariable("unitType")}
                    OnTitleUpdate={(title) => updateTitle(title)}
                />
            )
        }
        if (form === "unit_add") {
            return (
                <PageAddUnit
                    OnNavigate={(addr) => navigate(addr)}
                    UnitId={getHashVariable("unitId")}
                    OnTitleUpdate={(title) => updateTitle(title)}
                />
            )
        }
        if (form === "data_item") {
            return (
                <PageDataItem
                    OnNavigate={(addr) => navigate(addr)}
                    DataItemName={getHashVariable("dataItemName")}
                    OnTitleUpdate={(title) => updateTitle(title)}
                    FullWidth={mobileOpen}
                />
            )
        }
        if (form === "sensor_add") {
            return (
                <PageSensorAdd onComplete={() => {
                    navigate("#form=sensors")
                }}/>
            )
        }
        if (form === "account") {
            return (
                <PageAccount
                    OnNavigate={(addr) => navigate(addr)}
                    OnTitleUpdate={(title) => updateTitle(title)}
                    OnNeedUpdate={() => {
                        forceUpdate()
                    }}/>
            )
        }
        if (form === "about") {
            return (
                <PageAbout
                    OnNavigate={(addr) => navigate(addr)}
                    OnTitleUpdate={(title) => updateTitle(title)}
                    OnNeedUpdate={() => {
                        forceUpdate()
                    }}/>
            )
        }

        navigate("#form=units")

        return (
            <div>no form</div>
        )
    }

    const container = window !== undefined ? () => window().document.body : undefined;

    const th = createMuiTheme({
        palette: {
            type: 'dark',
            primary: {
                main: '#00A0E3'
            },
            secondary: {
                main: '#19BB4F'
            }
        }
    });

    if (getCookie("session_token") === undefined) {
        return (
            <MuiThemeProvider theme={th}>
                <div className={classes.root}>
                    <SignIn OnNeedUpdate={() => {
                        forceUpdate();
                    }}/>
                </div>
            </MuiThemeProvider>
        );
    }

    return (
        <MuiThemeProvider theme={th}>
            <SnackbarProvider maxSnack={3}>
                <div className={classes.root} key='gazer-main'>
                    <CssBaseline/>
                    <AppBar position="fixed" className={classes.appBar}
                            style={{backgroundColor: '#1E1E1E', color: '#CCC'}}>
                        <Toolbar>
                            <IconButton
                                color="inherit"
                                aria-label="open drawer"
                                edge="start"
                                onClick={handleDrawerToggle}
                                className={classes.menuButton}
                            >
                                <MenuIcon/>
                            </IconButton>
                            <Typography variant="h6" noWrap>
                                {title}
                            </Typography>
                        </Toolbar>
                    </AppBar>
                    <nav className={classes.drawer} aria-label="mailbox folders">
                        {/* The implementation can be swapped with js to avoid SEO duplication of links. */}
                        <Hidden smUp implementation="js">
                            <Drawer
                                container={container}
                                variant="temporary"
                                anchor={'left'}
                                open={mobileOpen}
                                onClose={handleDrawerToggle}
                                classes={{
                                    paper: classes.drawerPaper,
                                }}
                                ModalProps={{
                                    keepMounted: true, // Better open performance on mobile.
                                }}
                            >
                                {drawer}
                            </Drawer>
                        </Hidden>
                        <Hidden xsDown implementation="js">
                            <Drawer
                                classes={{
                                    paper: classes.drawerPaper,
                                }}
                                variant="permanent"
                                open
                            >
                                {drawer}
                            </Drawer>
                        </Hidden>
                    </nav>
                    <main className={classes.content}>
                        <div className={classes.toolbar}/>
                        {renderForm()}
                    </main>
                </div>
            </SnackbarProvider>
        </MuiThemeProvider>
    );
}

ResponsiveDrawer.propTypes = {
    /**
     * Injected by the documentation to work in an iframe.
     * You won't need it on your project.
     */
    window: PropTypes.func,
};

export default ResponsiveDrawer;
