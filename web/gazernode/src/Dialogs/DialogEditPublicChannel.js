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

export default function DialogEditPublicChannel(props){
    const [open, setOpen] = React.useState(false);
    const [error, setError] = React.useState("");
    const [channelName, setChannelName] = useState("")
    const {enqueueSnackbar} = useSnackbar();
    const [loading, setLoading] = React.useState(false)
    const [loaded, setLoaded] = React.useState(false)
    const [processing, setProcessing] = React.useState(false)

    const requestChannelState = () => {
        if (loading)
            return
        setLoading(true)

        let req = {
        }
        Request('public_channel_list', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            for (let i = 0; i < result.channels.length; i++)
                            {
                                let ch = result.channels[i]
                                if (ch.id === props.ChannelId) {
                                    setChannelName(ch.name)
                                    setLoaded(true)
                                }
                            }
                        }
                    );
                    setLoading(false)
                    return
                }
                if (res.status === 500) {
                    res.json().then(
                        (result) => {
                            enqueueSnackbar("Error: " + result.error, {variant: 'error'});
                        }
                    );
                    setLoading(false)
                    return
                }

                enqueueSnackbar("Error: " + res.status + ": " + res.statusText, {variant: 'error'});
                setLoading(false)
            }).catch((e)=> {
                enqueueSnackbar("Error: " + e.message(), {variant: 'error'});
        });
    }

    const requestEditChannel = () => {
        if (processing) {
            return
        }
        setProcessing(true)

        let req = {
            id: props.ChannelId,
            name: channelName,
        }
        Request('public_channel_set_name', req)
            .then((res) => {

                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setOpen(false)
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
        requestChannelState()
    };

    const handleClose = () => {
        setOpen(false);
    };

    const handleOK = () => {
        if (!loaded)
            return
        requestEditChannel()
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
                        disabled={!loaded}
                        fullWidth
                        placeholder="Channel name"
                        value={channelName}
                        onChange={(ev) => {
                            setChannelName(ev.target.value)
                        }}
                        onKeyDown={(e) => {
                            if (e.key === "Enter") {
                                handleOK()
                            }
                        }}
                    />
                    <LinearProgress style={{margin: "10px"}} hidden={!loading}/>
                    <LinearProgress style={{margin: "10px"}} hidden={!processing}/>
                </DialogContent>
                <DialogActions style={{backgroundColor: "#222"}}>
                    <Button onClick={handleOK} color="primary" variant="contained" style={{minWidth: "100px"}} disabled={!loaded}>
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
