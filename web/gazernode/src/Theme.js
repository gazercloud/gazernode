import { createMuiTheme } from '@material-ui/core/styles';
import green from "@material-ui/core/colors/green";
import purple from "@material-ui/core/colors/purple";

const theme = createMuiTheme({
    palette: {
        primary: green,
        secondary: purple,
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
