import React from 'react';
import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
import Link from '@material-ui/core/Link';
import Box from '@material-ui/core/Box';
import LockOutlinedIcon from '@material-ui/icons/LockOutlined';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';
import Request from "../request";

function Copyright() {
    return (
        <Typography variant="body2" color="textSecondary" align="center">
            {'Copyright © '}
            <Link color="inherit" href="https://gazer.cloud/">
                Gazer.Cloud
            </Link>{' '}
            {new Date().getFullYear()}
            {'.'}
        </Typography>
    );
}

const useStyles = makeStyles((theme) => ({
    paper: {
        marginTop: theme.spacing(8),
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
    },
    avatar: {
        margin: theme.spacing(1),
    },
    form: {
        width: '100%', // Fix IE 11 issue.
        marginTop: theme.spacing(1),
    },
    submit: {
        margin: theme.spacing(3, 0, 2),
    },
}));

export default function SignIn(props) {
    const classes = useStyles();
    const [userName, setUserName] = React.useState("")
    const [password, setPassword] = React.useState("")
    const [message, setMessage] = React.useState("")

    const requestLogin = () => {
        let req = {
            "user_name": userName,
            "password": password
        }
        Request('session_open', req)
            .then((res) => {
                if (res.status === 200) {
                    res.text().then(
                        (result) => {
                            try {
                                let obj = JSON.parse(result);
                                console.log("session_open ok", obj)
                                props.OnNeedUpdate()
                            } catch (e) {
                                setMessage("wrong server response")
                                console.log("session_open Wrong json", e)
                            }

                        }
                    )
                    return
                }
                if (res.status === 500) {
                    res.json().then(
                        (result) => {
                            setMessage("Error: " + result.error)
                            console.log("session_open ok", result)
                        }
                    );
                    return
                }

                res.text().then(
                    (result) => {
                        setMessage("Error " + res.status + ": " + result)
                    }
                )
            });

    }

    return (
        <Container component="main" maxWidth="xs">
            <CssBaseline />
            <div className={classes.paper}>
                <Avatar className={classes.avatar}>
                    <LockOutlinedIcon />
                </Avatar>
                <Typography component="h1" variant="h5">
                    Sign in
                </Typography>
                <form className={classes.form} noValidate>
                    <TextField
                        variant="outlined"
                        margin="normal"
                        required
                        fullWidth
                        id="email"
                        label="Email Address"
                        name="email"
                        autoComplete="email"
                        autoFocus
                        value={userName}
                        onChange={(event) => {
                            setUserName(event.target.value)
                        }}
                    />
                    <TextField
                        variant="outlined"
                        margin="normal"
                        required
                        fullWidth
                        name="password"
                        label="Password"
                        type="password"
                        id="password"
                        autoComplete="current-password"
                        value={password}
                        onChange={(event) => {
                            setPassword(event.target.value)
                        }}
                    />
                    <Button
                        fullWidth
                        variant="contained"
                        color="secondary"
                        className={classes.submit}
                        onClick={requestLogin}
                    >
                        Sign In
                    </Button>
                </form>
            </div>
            <div style={{color:'#F00', fontSize: '24pt'}}>
                {message}
            </div>
            <Box mt={8}>
                <Copyright />
            </Box>
        </Container>
    );
}
