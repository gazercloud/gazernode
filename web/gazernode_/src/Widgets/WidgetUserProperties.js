import React, {useEffect, useState} from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Paper from "@material-ui/core/Paper";
import Typography from "@material-ui/core/Typography";
import TextField from "@material-ui/core/TextField";
import Button from "@material-ui/core/Button";
import Grid from "@material-ui/core/Grid";
import Dialog from "@material-ui/core/Dialog";
import DialogTitle from "@material-ui/core/DialogTitle";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogActions from "@material-ui/core/DialogActions";
import Request from "../request";

const useStyles = makeStyles((theme) => ({
    root: {
        width: '100%',
        backgroundColor: theme.palette.background.paper,
    },
    values: {
        width: '100%',
        backgroundColor: theme.palette.background.paper,
    },
    expand: {
        transform: 'rotate(0deg)',
        marginLeft: 'auto',
        transition: theme.transitions.create('transform', {
            duration: theme.transitions.duration.shortest,
        }),
    },
    expandOpen: {
        transform: 'rotate(180deg)',
    },
}));



function WidgetUserProperties(props) {
    const classes = useStyles();
    const [userLoaded, setUserLoaded] = React.useState(false)

    const [userId, setUserId] = React.useState("")
    const [userName, setUserName] = React.useState("")
    const [userPassword, setUserPassword] = React.useState("")
    const [userDisplayName, setUserDisplayName] = React.useState("")
    const [userEmail, setUserEmail] = React.useState("")
    const [userPhone, setUserPhone] = React.useState("")

    const [userIdOriginal, setUserIdOriginal] = React.useState("")
    const [userNameOriginal, setUserNameOriginal] = React.useState("")
    const [userPasswordOriginal, setUserPasswordOriginal] = React.useState("")
    const [userDisplayNameOriginal, setUserDisplayNameOriginal] = React.useState("")
    const [userEmailOriginal, setUserEmailOriginal] = React.useState("")
    const [userPhoneOriginal, setUserPhoneOriginal] = React.useState("")

    const [errorMessage, setErrorMessage] = React.useState("")

    const [processingLoadUser, setProcessingLoadUser] = useState(false);
    const [needToLoadUser, setNeedToLoadUser] = useState(true);

    const [openConfirmation, setOpenConfirmation] = useState(false)

    const resetProperties = () => {
        setUserId("")

        setUserName("")
        setUserPassword("")
        setUserDisplayName("")
        setUserEmail("")
        setUserPhone("")

        setUserNameOriginal("")
        setUserPasswordOriginal("")
        setUserDisplayNameOriginal("")
        setUserEmailOriginal("")
        setUserPhoneOriginal("")
    }

    const reloadUser = () => {
        resetProperties()
        setNeedToLoadUser(true)
    }

    if (props.CurrentId !== userIdOriginal) {
        setUserIdOriginal(props.CurrentId)
        setUserLoaded(false)
        reloadUser()
        console.log("loading parameter")
        console.log(userIdOriginal)
    }

    const loadUser = () => {
        if (userIdOriginal === "")
            return

        if (!needToLoadUser)
            return
        if (processingLoadUser)
            return
        setUserLoaded(false)
        let req = {
            "func": "host_get_user",
            "id": userIdOriginal
        }
        setNeedToLoadUser(false)
        setProcessingLoadUser(true)

        Request('/api/request', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            let param = result
                            if (param.id === userIdOriginal) {

                                setUserId(param.id)
                                setUserName(param.name)
                                setUserPassword(param.password)
                                setUserDisplayName(param.display_name)
                                setUserEmail(param.email)
                                setUserPhone(param.phone)

                                setUserNameOriginal(param.name)
                                setUserPasswordOriginal(param.password)
                                setUserDisplayNameOriginal(param.display_name)
                                setUserEmailOriginal(param.email)
                                setUserPhoneOriginal(param.phone)

                                setUserLoaded(true)
                                setErrorMessage("")
                            }
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            setErrorMessage(result.error)
                        }
                    );
                }
                setProcessingLoadUser(false)
            });
    }

    const denyEdit = () => {
        return !userLoaded
    }

    const denySave = () => {
        if (processingLoadUser)
            return true
        if (!userLoaded)
            return true
        if (!hasChanges())
            return true
        return false
    }

    const updateUser = () => {
        let req = {
            "func": "host_update_user",
            "id": userIdOriginal,
            "name": userName,
            "password": userPassword,
            "display_name": userDisplayName,
            "email": userEmail,
            "phone": userPhone
        }
        console.log(req)

        Request('/api/request', req)
            .then((res) => {
                if (res.status === 200) {
                    reloadUser();
                } else {
                    res.json().then(
                        (result) => {
                            setErrorMessage(result.error)
                        }
                    );
                }
            });
    }

    const removeUser = () => {
        let req = {
            "func": "host_remove_user",
            "id": userIdOriginal
        }

        Request('/api/request', req)
            .then((res) => {
                if (res.status === 200) {
                    props.OnNeedToLoadUserList()
                } else {
                    res.json().then(
                        (result) => {
                            setErrorMessage(result.error)
                        }
                    );
                }
            });
        setOpenConfirmation(false)
    }

    useEffect(() => {
        const timer = setInterval(() => {
            loadUser()
        }, 100);
        return () => clearInterval(timer);
    });

    const hasChanges = () => {
        if (userId !== userIdOriginal)
            return true
        if (userName !== userNameOriginal)
            return true
        if (userPassword !== userPasswordOriginal)
            return true
        if (userDisplayName !== userDisplayNameOriginal)
            return true
        if (userEmail !== userEmailOriginal)
            return true
        if (userPhone !== userPhoneOriginal)
            return true
        return false
    }

    return (
        <Paper className={classes.root} variant="outlined"
               style={{margin: "10px", padding: "5px", minWidth: "360px"}}>
            <Grid container direction="column">
                <Grid container direction="column">
                    <Grid item>
                        <TextField
                            disabled={denyEdit()}
                            style={{margin: "10px", width: "250px"}}
                            label="Name"
                            value={userName}
                            variant="outlined"
                            onChange={(event) => {
                                setUserName(event.target.value)
                            }}>
                        </TextField></Grid>
                    <Grid item>
                        <Grid container direction="row">
                            <Grid item>
                                <TextField
                                    disabled={denyEdit()}
                                    style={{margin: "10px", width: "250px"}}
                                    label="Password"
                                    value={userPassword}
                                    variant="outlined"
                                    onChange={(event) => {
                                        setUserPassword(event.target.value)
                                    }}
                                /></Grid>
                            <Grid item style={{margin: "10px", fontStyle: "italic", fontSize: "10pt"}}>
                            </Grid>
                        </Grid>
                    </Grid>
                    <Grid item>
                        <Grid container direction="row">
                            <Grid item>
                                <TextField
                                    disabled={denyEdit()}
                                    style={{margin: "10px", width: "250px"}}
                                    label="DisplayName"
                                    value={userDisplayName}
                                    variant="outlined"
                                    onChange={(event) => {
                                        setUserDisplayName(event.target.value)
                                    }}/>
                            </Grid>
                            <Grid item style={{margin: "10px", fontStyle: "italic", fontSize: "10pt"}}>
                            </Grid>
                        </Grid>
                    </Grid>
                    <Grid item>
                        <TextField disabled={denyEdit()}
                                   style={{margin: "10px", width: "250px"}}
                                   label="Email"
                                   value={userEmail}
                                   variant="outlined" onChange={(event) => {
                            setUserEmail(event.target.value)
                        }}/>
                    </Grid>
                    <Grid item>
                        <TextField disabled={denyEdit()}
                                   style={{margin: "10px", width: "250px"}}
                                   label="Phone"
                                   value={userPhone}
                                   variant="outlined" onChange={(event) => {
                            setUserPhone(event.target.value)
                        }}/></Grid>
                </Grid>
                <Grid item>
                    <Button
                        variant="contained"
                        color="primary"
                        disabled={denySave()}
                        style={{margin: "10px", width: "250px"}}
                        onClick={updateUser}>
                        Save
                    </Button>
                    <Button
                        variant="contained"
                        color="primary"
                        disabled={denyEdit()}
                        style={{margin: "10px", width: "100px"}}
                        onClick={() => { setOpenConfirmation(true) }}>
                        Remove
                    </Button>
                </Grid>
                <Grid item>
                    <Typography hidden={errorMessage.length < 1} style={{color: "#FF2906", margin: "10px"}}>Error: <br/>{errorMessage}</Typography>
                </Grid>
            </Grid>
            <Dialog open={openConfirmation} onClose={() => { setOpenConfirmation(false) }} aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">Remove user confirmation</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Remove user: {userNameOriginal} ?
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button variant="outlined" onClick={removeUser} color="primary">
                        OK
                    </Button>
                    <Button onClick={() => { setOpenConfirmation(false) }} color="primary">
                        Cancel
                    </Button>
                </DialogActions>
            </Dialog>
        </Paper>
    );
}

export default WidgetUserProperties;

