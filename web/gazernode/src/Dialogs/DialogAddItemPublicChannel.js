import React, {useState} from "react";
import Request from "../request";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableRow from "@material-ui/core/TableRow";
import {Tooltip} from "@material-ui/core";
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

export default function DialogAddItemPublicChannel(props){
    const [open, setOpen] = React.useState(false);
    const [error, setError] = React.useState("");
    const [currentItem, setCurrentItem] = useState("")
    const [currentChannelId, setCurrentChannelId] = useState("")
    const [hoverItem, setHoverItem] = useState("")
    const [channels, setChannels] = React.useState({})

    const requestChannels = () => {
        let req = {
        }
        Request('public_channel_list', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setChannels(result)
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                }
            });
    }

    const btnStyle = (key) => {
        if (currentItem === key) {
            return {
                cursor: "pointer",
                backgroundColor: "#52bdff",
            }
        } else {
            if (hoverItem === key) {
                return {
                    cursor: "pointer",
                    backgroundColor: "#b8e4ff"
                }
            } else {
                return {
                    cursor: "pointer",
                    backgroundColor: "#222"
                }
            }
        }
    }

    const requestAddItem = () => {
        let req = {
            ids: [currentChannelId],
            items: [props.ItemName]
        }
        Request('public_channel_item_add', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            props.OnSuccess(currentChannelId)
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                }
            });
    }

    const handleOpen = () => {
        setOpen(true);
        requestChannels();
    };

    const handleClose = () => {
        setOpen(false);
    };

    const handleOK = () => {
        requestAddItem(currentChannelId)
        setOpen(false)
    };

    const handleEnter = (ev, key) => {
        setHoverItem(ev)
    }

    const handleLeave = (ev, key) => {
        setHoverItem("")
    }

    const btnClick = (ev, key) => {
        setCurrentItem(ev)
        setCurrentChannelId(key)
    }

    const drawContent = () => {
        if (channels === undefined) {
            return (
                <div>no data ... loading ...</div>
            )
        }

        return (
            <div>
                <Table size="small" style={{minWidth: "300px"}}>
                    <TableBody>
                        {(channels.channels !== undefined) ? channels.channels.map((dataItem, dataItemIndex) =>
                            (
                                <TableRow
                                    onMouseEnter={() => handleEnter(dataItemIndex)}
                                    onMouseLeave={() => handleLeave(dataItemIndex)}
                                    style={btnStyle(dataItemIndex)}
                                    onClick={btnClick.bind(this, dataItemIndex, dataItem.id)}
                                >
                                    <Typography style={{padding: "5px"}}>{dataItem.name}</Typography>
                                </TableRow>
                            )) : <div/>
                        }
                    </TableBody>
                </Table>

            </div>
        )
    }

    if (props.Entity === "") {
        return (
            <div></div>
        )
    }

    return (
        <span>
            <Tooltip title="Add to public channel" TransitionComponent={Zoom} style={{border: "1px solid #333", marginRight: "5px"}}>
                <IconButton onClick={handleOpen} variant="outlined" color="primary"><BackupOutlinedIcon  fontSize="large"/></IconButton>
            </Tooltip>
            <Dialog open={open} onClose={handleClose} aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">Add to public channel</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Lookup: {props.Entity}
                    </DialogContentText>
                    <DialogContent>
                        {drawContent()}
                    </DialogContent>
                </DialogContent>
                <DialogActions>
                    <Button onClick={handleClose} color="primary">
                        Cancel
                    </Button>
                    <Button onClick={handleOK} color="primary">
                        OK
                    </Button>
                </DialogActions>
            </Dialog>

            <div>{error}</div>
        </span>
    );
}
