import React, {useState} from "react";
import Request from "../request";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableRow from "@material-ui/core/TableRow";
import {LinearProgress, TextField, Tooltip} from "@material-ui/core";
import Zoom from "@material-ui/core/Zoom";
import Button from "@material-ui/core/Button";
import Dialog from "@material-ui/core/Dialog";
import DialogTitle from "@material-ui/core/DialogTitle";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogActions from "@material-ui/core/DialogActions";
import BlurOnIcon from "@material-ui/icons/BlurOn";
import IconButton from "@material-ui/core/IconButton";
import BackupOutlinedIcon from '@material-ui/icons/BackupOutlined';
import Typography from "@material-ui/core/Typography";
import AddOutlinedIcon from '@material-ui/icons/AddOutlined';
import EditIcon from '@material-ui/icons/Edit';
import {useSnackbar} from "notistack";
import Grid from "@material-ui/core/Grid";
import Hidden from "@material-ui/core/Hidden";

export default function DialogEditUser(props){
    const [open, setOpen] = React.useState(false);
    const [error, setError] = React.useState("");
    const [password1, setPassword1] = useState("")
    const [password2, setPassword2] = useState("")
    const {enqueueSnackbar} = useSnackbar();
    const [processing, setProcessing] = React.useState(false)

    const requestEditUser = () => {
        if (password1 !== password2) {
            enqueueSnackbar("Error: Password mismatch", {variant: 'error'});
            return
        }

        if (processing) {
            return
        }
        setProcessing(true)

        let req = {
            user_name: props.UserName,
            password: password1,
        }
        Request('user_set_password', req)
            .then((res) => {

                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setOpen(false)
                            enqueueSnackbar("Password changed", {variant: 'success'});
                        }
                    );
                    setProcessing(false)
                    return
                }

                if (res.status === 500) {
                    res.json().then(
                        (result) => {
                            enqueueSnackbar("Error: " + result.error, {variant: 'error'});
                        }
                    );
                    setProcessing(false)
                    return
                }

                enqueueSnackbar("Error: " + res.status + ": " + res.statusText, {variant: 'error'});
                setProcessing(false)
            });
    }

    const handleOpen = () => {
        setOpen(true);
    };

    const handleClose = () => {
        setOpen(false);
    };

    const handleOK = () => {
        requestEditUser()
    };

    return (
        <span>
            <Tooltip title="Edit public channel" TransitionComponent={Zoom} style={{border: "1px solid #333", marginRight: "5px"}}>
                <IconButton onClick={handleOpen} variant="outlined" color="primary">
                    <EditIcon  fontSize="large"/>
                </IconButton>
            </Tooltip>
            <Dialog open={open} onClose={handleClose} aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">
                                        <Grid container alignItems="center">
                        <EditIcon fontSize="large" style={{marginRight: "10px"}}/>
                        <span  style={{margin: "0px"}}>Edit Public Channel</span>
                    </Grid>
                </DialogTitle>
                <DialogContent style={{minWidth: "350px", backgroundColor: "#222", padding: "20px"}}>
                    <TextField
                        fullWidth
                        type="password"
                        placeholder="Password"
                        value={password1}
                        onChange={(ev) => {
                            setPassword1(ev.target.value)
                        }}
                        onKeyDown={(e) => {
                            if (e.key === "Enter") {
                                handleOK()
                            }
                        }}
                    />
                    <TextField
                        fullWidth
                        type="password"
                        placeholder="Retype Password"
                        value={password2}
                        onChange={(ev) => {
                            setPassword2(ev.target.value)
                        }}
                        onKeyDown={(e) => {
                            if (e.key === "Enter") {
                                handleOK()
                            }
                        }}
                    />
                    <LinearProgress style={{margin: "10px"}} hidden={!processing}/>
                </DialogContent>
                <DialogActions style={{backgroundColor: "#222"}}>
                    <Button onClick={handleOK} color="primary" variant="contained" style={{minWidth: "100px"}}>
                        OK
                    </Button>
                    <Button onClick={handleClose} color="primary" variant="outlined" style={{minWidth: "100px"}}>
                        Cancel
                    </Button>
                </DialogActions>
            </Dialog>

            <span>{error}</span>
        </span>
    );
}
