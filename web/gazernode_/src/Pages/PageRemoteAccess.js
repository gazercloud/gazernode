import React, {useEffect, useState} from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Grid from "@material-ui/core/Grid";
import Request, {RequestFailed} from "../request";
import DialogAddPublicChannel from "../Dialogs/DialogAddPublicChannel";
import {useSnackbar} from "notistack";
import {Button, LinearProgress, Paper, TableContainer, TextField, Typography} from "@material-ui/core";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableRow from "@material-ui/core/TableRow";
import TableCell from "@material-ui/core/TableCell";
import WidgetError from "../Widgets/WidgetError";

export default function PageRemoteAccess(props) {
    const [cloudState, setCloudState] = React.useState({})
    const { enqueueSnackbar } = useSnackbar();
    const [userName, setUserName] = React.useState("")
    const [password, setPassword] = React.useState("")
    const [stateLoaded, setStateLoaded] = React.useState(false)
    const [loginProcessing, setLoginProcessing] = React.useState(false)
    const [logoutProcessing, setLogoutProcessing] = React.useState(false)
    const [messageError, setMessageError] = React.useState("")

    const requestCloudState = () => {
        let req = {
        }
        Request('cloud_state', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setCloudState(result)
                            setStateLoaded(true)
                            setMessageError("")
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            setMessageError(result.error)
                        }
                    );
                }
            }).catch(res => {
            setMessageError(res.message)
            RequestFailed()
        });
    }

    const requestCloudLogin = () => {
        if (loginProcessing) {
            return
        }
        setLoginProcessing(true)
        let req = {
            user_name: userName,
            password: password
        }
        Request('cloud_login', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            requestCloudState()
                            setLoginProcessing(false)
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                            setLoginProcessing(false)
                        }
                    );
                }
                setLoginProcessing(false)
            });
    }

    const requestCloudLogout = () => {
        if (logoutProcessing) {
            return
        }
        setLogoutProcessing(true)

        let req = {
        }
        Request('cloud_logout', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            requestCloudState()
                            setLogoutProcessing(false)
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                            setLogoutProcessing(false)
                            setLogoutProcessing(false)
                        }
                    );
                }
            });
    }

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        props.OnTitleUpdate("Remove Access")
        requestCloudState()
        setFirstRendering(false)
    }

    useEffect(() => {
        const timer = setInterval(() => {
            requestCloudState()
        }, 1000);
        return () => clearInterval(timer);
    });

    const nodeStatus = () => {
        if (cloudState === undefined) {
            return (
                <div style={{color: "#F40"}}>"unknown"</div>
            )
        }
        if (cloudState.i_am_status === "ok") {
            return (
                <div style={{color: "#080"}}>{cloudState.node_id} / ok</div>
            )
        }
        return (
            <div style={{color: "#F40"}}>{cloudState.node_id} / {cloudState.i_am_status}</div>
        )
    }

    const connectedStatus = () => {
        if (cloudState === undefined) {
            return (
                <div style={{color: "#F40"}}>unknown</div>
            )
        }
        if (cloudState.connected) {
            return (
                <div style={{color: "#080"}}>{cloudState.current_repeater} / ok</div>
            )
        }
        return (
            <div style={{color: "#F40"}}>{cloudState.current_repeater} / {cloudState.connection_status}</div>
        )
    }

    const loginStatus = () => {
        if (cloudState === undefined) {
            return (
                <div style={{color: "#F40"}}>unknown</div>
            )
        }
        if (cloudState.logged_in) {
            return (
                <div style={{color: "#080"}}>{cloudState.user_name} / ok</div>
            )
        }
        return (
            <div style={{color: "#F40"}}>{cloudState.user_name} / {cloudState.login_status}</div>
        )
    }

    const displayStatus = () => {
        return (
            <div>
                <Grid container direction="row">
                    <Grid item>
                        <TableContainer style={{maxWidth: "600px", maxHeight: "400px"}}>
                            <Table size="small" style={{minWidth: "400px"}}>

                                <TableBody>
                                    <TableRow>
                                        <TableCell style={{maxWidth: "200px", fontSize: "20pt"}}>
                                            Node Id
                                        </TableCell>
                                        <TableCell style={{maxWidth: "500px", fontSize: "20pt"}}>
                                            {nodeStatus()}
                                        </TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell style={{maxWidth: "200px", fontSize: "20pt"}}>
                                            User Name
                                        </TableCell>
                                        <TableCell style={{maxWidth: "500px", fontSize: "20pt"}}>
                                            {loginStatus()}
                                        </TableCell>
                                    </TableRow>
                                    <TableRow>
                                        <TableCell style={{maxWidth: "200px", fontSize: "20pt"}}>
                                            Connection
                                        </TableCell>
                                        <TableCell style={{maxWidth: "500px", fontSize: "20pt"}}>
                                            {connectedStatus()}
                                        </TableCell>
                                    </TableRow>
                                </TableBody>
                            </Table>
                        </TableContainer>
                    </Grid>
                </Grid>
            </div>
        );
    }

    const displayLoginForm = () => {
        if (cloudState === undefined) {
            return (
                <div/>
            )
        }

        if (!cloudState.logged_in) {
            return (
                <div style={{marginTop: "20px"}}>
                    <Typography style={{fontSize: "24pt", color: "#F80"}}>Login to the Gazer Cloud</Typography>
                    <Grid container direction="column" style={{
                        maxWidth: "400px",
                        backgroundColor: "#333",
                        padding: "10px",
                        borderRadius: "10px",
                        marginTop: "10px"
                    }}>
                        <TextField fullWidth placeholder="E-Mail" style={{padding: "10px"}} value={userName}
                                   onChange={(e) => {
                                       setUserName(e.target.value)
                                   }}/>
                        <TextField fullWidth type="password" placeholder="Password" style={{padding: "10px"}}
                                   value={password} onChange={(e) => {
                            setPassword(e.target.value)
                        }}/>
                        <Button variant="outlined" style={{margin: "10px"}} onClick={requestCloudLogin}>Login</Button>
                        <LinearProgress hidden={!loginProcessing}/>
                    </Grid>
                </div>
            )
        } else {
            return (
                <div>
                    <Button variant="outlined" style={{margin: "10px"}} onClick={requestCloudLogout} color="primary">
                        Logout
                    </Button>
                    <LinearProgress hidden={!logoutProcessing}/>

                </div>
            )
        }
    }

    if (stateLoaded) {
        return (
            <div>
                {displayStatus()}
                {displayLoginForm()}
                <div>
                    <WidgetError Message={messageError} />
                </div>
            </div>
        );
    }

    return (
        <div>
            <LinearProgress/>
            <div>
                <WidgetError Message={messageError} />
            </div>
        </div>
    );
}
