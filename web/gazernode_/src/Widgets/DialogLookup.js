import React, {useState} from 'react';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';
import Request from "../request";
import IconButton from "@material-ui/core/IconButton";
import EditIcon from '@material-ui/icons/Edit';
import Zoom from "@material-ui/core/Zoom";
import {Tooltip} from "@material-ui/core";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableRow from "@material-ui/core/TableRow";
import TableCell from "@material-ui/core/TableCell";
import MoreHorizIcon from '@material-ui/icons/MoreHoriz';

export default function DialogLookup(props) {
    const [open, setOpen] = React.useState(false);
    const [text, setText] = React.useState("");
    const [selectedItem, setSelectedItem] = React.useState("");
    const [error, setError] = React.useState("");
    const [lookupResult, setLookupResult] = React.useState(undefined)
    const [currentItem, setCurrentItem] = useState("")
    const [hoverItem, setHoverItem] = useState("")

    const requestItems = () => {
        let req = {
            entity: props.Entity,
            parameters: props.Parameters
        }
        Request('service_lookup', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setLookupResult(result.result)
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            setError(result.error)
                        }
                    );
                }
            }).catch((e) => {
            setError(e.message)
        });
    }

    const handleOpen = () => {
        setOpen(true);
        requestItems();
    };

    const handleClose = () => {
        setOpen(false);
    };

    const handleOK = () => {
        let keyColumnIndex = 0
        for (let i = 0; i < lookupResult.columns.length; i++)
        {
            if (lookupResult.columns[i] === lookupResult.key_column) {
                keyColumnIndex = i;
                break;
            }
        }

        props.OnSelected(lookupResult.rows[currentItem].cells[keyColumnIndex])
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

    const drawContent = () => {
        if (lookupResult === undefined) {
            return (
                <div>no data ... loading ...</div>
            )
        }

        return (
            <div>
                <div>{lookupResult.entity}</div>
                <Table size="small" style={{minWidth: "300px"}}>
                    <TableBody>
                        <TableRow>
                            {lookupResult.columns.map((item) => (
                                <TableCell>
                                    {lookupResult.key_column === item.name ?
                                        <div style={{fontWeight: "bold"}}>{item.display_name}</div>
                                        :
                                        <div>{item.display_name}</div>
                                    }
                                </TableCell>
                            ))
                            }
                        </TableRow>
                        {lookupResult.rows !== undefined ? lookupResult.rows.map((dataItem, dataItemIndex) =>
                            (
                                <TableRow
                                    onMouseEnter={() => handleEnter(dataItemIndex)}
                                    onMouseLeave={() => handleLeave(dataItemIndex)}
                                    style={btnStyle(dataItemIndex)}
                                    onClick={btnClick.bind(this, dataItemIndex)}
                                >
                                    {dataItem.cells.map((item) => (
                                        <TableCell>
                                            <div>{item}</div>
                                        </TableCell>
                                    ))}
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
        <div>
            <Tooltip title="Select" TransitionComponent={Zoom}><IconButton><MoreHorizIcon fontSize="small" onClick={handleOpen} /></IconButton></Tooltip>
            <Dialog open={open} onClose={handleClose} aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">Lookup</DialogTitle>
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
        </div>
    );
}
