import React, {useState} from "react";
import Request from "../request";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableRow from "@material-ui/core/TableRow";
import {TextField, Tooltip} from "@material-ui/core";
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
import {SaveAltOutlined} from "@material-ui/icons";

export default function DialogWriteValue(props){
    const [open, setOpen] = React.useState(false);
    const [error, setError] = React.useState("");
    const [currentItem, setCurrentItem] = useState("")
    const [currentChannelId, setCurrentChannelId] = useState("")
    const [hoverItem, setHoverItem] = useState("")
    const [channels, setChannels] = React.useState({})
    const [value, setValue] = React.useState("")

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

    const requestWriteValue = () => {
        let req = {
            item_name: props.ItemName,
            value: value
        }
        Request('data_item_write', req)
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
    };

    const handleClose = () => {
        setOpen(false);
    };

    const handleOK = () => {
        requestWriteValue()
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
        return (
            <div>
                <TextField
                    value={value} onChange={(ev)=> {
                    setValue(ev.target.value)
                }}
                    autoFocus={true}
                    placeholder="Value"
                    onKeyUp={(ev) => {
                        if (ev.key === 'Enter') {
                            ev.preventDefault();
                            handleOK();
                        }
                    }}
                    onFocus={event => {
                        event.target.select();
                    }}
                />
            </div>
        )
    }

    return (
        <span>
            <Tooltip title="Write Value" TransitionComponent={Zoom} style={{border: "1px solid #333", marginRight: "5px"}}>
                <IconButton onClick={handleOpen} variant="outlined" color="primary"><SaveAltOutlined  fontSize="large"/></IconButton>
            </Tooltip>
            <Dialog open={open} onClose={handleClose} aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">Write Value</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Write Value: {props.ItemName}
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
