import React from 'react';
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
import { makeStyles, useTheme } from '@material-ui/core/styles';
import AccountTreeIcon from '@material-ui/icons/AccountTree';
import LayersIcon from '@material-ui/icons/Layers';
import TimelineIcon from '@material-ui/icons/Timeline';
import BlurOnIcon from '@material-ui/icons/BlurOn';
import PeopleIcon from '@material-ui/icons/People';
import AllInclusiveIcon from '@material-ui/icons/AllInclusive';
import InfoOutlinedIcon from '@material-ui/icons/InfoOutlined';
import PageSensors from "./Pages/PageSensors";
import PageAbout from "./Pages/PageAbout";
import PageDataItems from "./Pages/PageDataItems";
import PageMaps from "./Pages/PageMaps";
import PageCharts from "./Pages/PageCharts";
import PageUsers from "./Pages/PageUsers";
import PageAdmin from "./Pages/PageAdmin";
import Grid from "@material-ui/core/Grid";
import PageSensorAdd from "./Pages/PageSensorAdd";
import PageRequest from "./Pages/PageRequest";

const drawerWidth = 240;

const useStyles = makeStyles((theme) => ({
    root: {
        display: 'flex',
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
            backgroundColor: "#EEEEEE",
            color: "#000000"
        },
    },
    menuButton: {
        marginRight: theme.spacing(2),
        [theme.breakpoints.up('sm')]: {
            display: 'none',
        },
    },
    // necessary for content to be below app bar
    toolbar: theme.mixins.toolbar,
    drawerPaper: {
        width: drawerWidth,
        color: "#EEEEEE",
        backgroundColor: "#374859"
    },
    content: {
        flexGrow: 1,
        padding: theme.spacing(3),
    },
}));

function getHashVariable(variable) {
    const query = window.location.hash.substring(1);
    const vars = query.split('&');
    console.log(vars)
    for (let i = 0; i < vars.length; i++) {
        const pair = vars[i].split('=');
        console.log(pair[0])
        console.log(pair[1])
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

function ResponsiveDrawer(props) {
    const { window } = props;
    const classes = useStyles();
    const theme = useTheme();
    const [mobileOpen, setMobileOpen] = React.useState(false);
    const [redrawBool, setRedrawBool] = React.useState(false);

    const navigate = (link) => {
        gotoLink(link)
        setRedrawBool(!redrawBool)
    }

    const handleDrawerToggle = () => {
        setMobileOpen(!mobileOpen);
    };

    const drawer = (
        <div>
            <Grid container className={classes.toolbar} direction="row" alignItems="center" alignContent="center" >

                <Grid container direction="row" alignItems="flex-start" alignContent="flex-start" style={{margin: "15px"}}>
                    <Grid item><a href="/"><img src="/mainicon32.png" style={{width: "48px"}}/></a></Grid>
                    <Grid item>
                        <Typography style={{marginLeft: "8px", marginBottom: "0px", fontSize: "14pt"}}>
                            allece
                        </Typography>
                        <Typography style={{marginLeft: "8px", fontSize: "10pt", color: "#888888"}}>
                            watch your data
                        </Typography>
                    </Grid>
                </Grid>

            </Grid>
            <Divider />
            <List>
                <ListItem button key="sensors" component="a" onClick={() => {navigate("#form=sensors")}}>
                    <ListItemIcon><BlurOnIcon style={{color:"#FFFFFF"}}/></ListItemIcon>
                    <ListItemText primary={"Sensors"} />
                </ListItem>
                <ListItem button key="data_items" component="a" onClick={() => {navigate("#form=data_items")}}>
                    <ListItemIcon><AccountTreeIcon style={{color:"#FFFFFF"}} /></ListItemIcon>
                    <ListItemText primary={"Data Items"} />
                </ListItem>
                <ListItem button key="maps" component="a" onClick={() => {navigate("#form=maps")}}>
                    <ListItemIcon><LayersIcon style={{color:"#FFFFFF"}}/></ListItemIcon>
                    <ListItemText primary={"Maps"} />
                </ListItem>
                <ListItem button key="charts" component="a" onClick={() => {navigate("#form=charts")}}>
                    <ListItemIcon><TimelineIcon style={{color:"#FFFFFF"}}/></ListItemIcon>
                    <ListItemText primary={"Charts"} />
                </ListItem>
            </List>
            <Divider />
            <List>
                <ListItem button key="users" component="a" onClick={() => {navigate("#form=users")}}>
                    <ListItemIcon><PeopleIcon style={{color:"#FFFFFF"}}/></ListItemIcon>
                    <ListItemText primary={"Users"} />
                </ListItem>
                <ListItem button key="admin" component="a" onClick={() => {navigate("#form=admin")}}>
                    <ListItemIcon><AllInclusiveIcon style={{color:"#FFFFFF"}}/></ListItemIcon>
                    <ListItemText primary={"Administration"} />
                </ListItem>
                <ListItem button key="request" component="a" onClick={() => {navigate("#form=request")}}>
                    <ListItemIcon><AllInclusiveIcon style={{color:"#FFFFFF"}}/></ListItemIcon>
                    <ListItemText primary={"Request"} />
                </ListItem>
                <ListItem button key="about" component="a" onClick={() => {navigate("#form=about")}}>
                    <ListItemIcon><InfoOutlinedIcon style={{color:"#FFFFFF"}}/></ListItemIcon>
                    <ListItemText primary={"About"} />
                </ListItem>
            </List>
        </div>
    );




    const renderForm = () => {
        const form = getHashVariable("form")
        if (form === "sensors") {
            return (
                <PageSensors onAddSensor={() => {navigate("#form=sensor_add")}} />
            )
        }
        if (form === "sensor_add") {
            return (
                <PageSensorAdd onComplete={() => {navigate("#form=sensors")}} />
            )
        }
        if (form === "data_items") {
            return (
                <PageDataItems />
            )
        }
        if (form === "maps") {
            return (
                <PageMaps />
            )
        }
        if (form === "charts") {
            return (
                <PageCharts />
            )
        }

        if (form === "users") {
            return (
                <PageUsers />
            )
        }

        if (form === "admin") {
            return (
                <PageAdmin />
            )
        }

        if (form === "request") {
            return (
                <PageRequest />
            )
        }

        if (form === "about") {
            return (
                <PageAbout />
            )
        }
        return (
            <div>no form</div>
        )
    }

    const container = window !== undefined ? () => window().document.body : undefined;

    return (
        <div className={classes.root}>
            <CssBaseline />
            <AppBar position="fixed" className={classes.appBar}>
                <Toolbar>
                    <IconButton
                        color="inherit"
                        aria-label="open drawer"
                        edge="start"
                        onClick={handleDrawerToggle}
                        className={classes.menuButton}
                    >
                        <MenuIcon />
                    </IconButton>
                    <Typography variant="h6" noWrap>
                        Allece
                    </Typography>
                </Toolbar>
            </AppBar>
            <nav className={classes.drawer} aria-label="mailbox folders">
                {/* The implementation can be swapped with js to avoid SEO duplication of links. */}
                <Hidden smUp implementation="css">
                    <Drawer
                        container={container}
                        variant="temporary"
                        anchor={theme.direction === 'rtl' ? 'right' : 'left'}
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
                <Hidden xsDown implementation="css">
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
                <div className={classes.toolbar} />
                {renderForm()}
            </main>
        </div>
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
