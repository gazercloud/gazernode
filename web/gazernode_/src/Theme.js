import { createMuiTheme } from '@material-ui/core/styles';
import green from "@material-ui/core/colors/green";
import purple from "@material-ui/core/colors/purple";

const theme = createMuiTheme({
    palette: {
        primary: {
            main: '#E30000'
        },
        secondary: {
            main: '#00FF7F'
        }
    },

    status: {
        danger: 'orange',
    },
    overrides: {
    MuiDrawer: {
        background: '#FF202c',
    },
},
});

export default theme;
