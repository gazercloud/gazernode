import React, {useEffect, useState} from 'react';
import Request, {RequestFailed} from "../request";
import WidgetItemHistory from "../Widgets/WidgetItemHistory";
import DialogAddItemPublicChannel from "../Dialogs/DialogAddItemPublicChannel";
import {Tooltip} from "@material-ui/core";
import Zoom from "@material-ui/core/Zoom";
import IconButton from "@material-ui/core/IconButton";
import ArrowBackIosOutlinedIcon from '@material-ui/icons/ArrowBackIosOutlined';
import Grid from "@material-ui/core/Grid";
import Divider from "@material-ui/core/Divider";
import {useSnackbar} from "notistack";
import WidgetError from "../Widgets/WidgetError";
import DialogWriteValue from "../Dialogs/DialogWriteValue";
import {HotKeys} from "react-hotkeys";


export default function PageDataItem(props) {
    const [dataItemState, setDataItemState] = React.useState([])
    const [messageError, setMessageError] = React.useState("")
    const [stateLoading, setStateLoading] = React.useState(false)

    const requestDataItem = (dataItemName) => {
        if (stateLoading)
            return

        setStateLoading(true)

        let req = {
            items: [dataItemName]
        }
        Request('data_item_list', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            if (result.items.length > 0) {
                                setDataItemState(result.items[0])
                            }
                            setMessageError("")
                            setStateLoading(false)
                        }
                    );
                    return
                }

                if (res.status === 500) {
                    res.json().then(
                        (result) => {
                            setMessageError(result.error)
                            setStateLoading(false)
                        }
                    );
                    RequestFailed()
                    setStateLoading(false)
                    return
                }

                setMessageError("Error: " + res.status + " " + res.statusText)
                RequestFailed()
                setStateLoading(false)

            }).catch(res => {
            setMessageError(res.message)
            setStateLoading(false)
            RequestFailed()
        });
    }

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        const dataItemName = new Buffer(props.DataItemName, 'hex').toString();
        requestDataItem(dataItemName)
        props.OnTitleUpdate("Data Item [" + dataItemName + "]")
        setFirstRendering(false)
    }

    const getDataItemName = () => {
        if (props.DataItemName === undefined || props.DataItemName === "")
            return ""
        return new Buffer(props.DataItemName, 'hex').toString();
    }

    useEffect(() => {
        const dataItemName = getDataItemName();
        const timer = setInterval(() => {
            requestDataItem(dataItemName)
        }, 1000);
        return () => clearInterval(timer);
    });

    const formatTime = (date) => {

        let hh = date.getHours();
        if (hh < 10) hh = '0' + hh;

        let mm = date.getMinutes();
        if (mm < 10) mm = '0' + mm;

        let ss = date.getSeconds();
        if (ss < 10) ss = '0' + ss;

        let fff = date.getMilliseconds();
        if (fff < 10) {
            fff = '00' + fff;
        } else {
            if (fff < 100)
                fff = '0' + fff
        }

        return hh + ':' + mm + ':' + ss + '.' + fff;
    }

    const displayItemValue = (item) => {
        let dt1 = new Date(item.value.t / 1000)
        let dt2 = ""
        if (dt1 !== undefined) {
            dt2 = formatTime(dt1)
        }

        if (item.value.u === "error") {
            return (
                <div>
                    <div style={{
                        color: '#F30',

                        fontSize: '36pt',
                        textAlign: 'left'
                    }}>
                        {item.value.v + " " + item.value.u}
                    </div>
                    <div>
                        <div style={{fontSize: '14pt', color: '#AAA'}}>
                            <span>{dt2}</span>
                        </div>
                    </div>
                </div>
            )
        }

        return (
            <div>
                <div style={{
                    color: '#080',

                    fontSize: '36pt',
                    textAlign: 'left'
                }}>
                    {item.value.v + "  " + item.value.u}
                </div>
                <div>
                    <div style={{fontSize: '14pt', color: '#AAA'}}>
                        <span>{dt2}</span>
                    </div>
                </div>
            </div>
        )
    }

    const displayItemName = (item) => {
        return (
            <div style={{fontSize: '20pt'}}>{item.name}</div>
        )
    }

    const displayItem = (item) => {
        return (
            <div>
                {displayItemName(item)}
                {displayItemValue(item)}
            </div>
        )
    }

    const displayItems = () => {

        if (dataItemState === undefined || dataItemState.value === undefined) {
            return (
                <div>
                    loading ...
                </div>
            )
        }

        return (
            <div>
                <div style={{
                    padding: '20px',
                    borderRadius: '10px',
                    marginTop: '10px',
                    marginBottom: '10px',
                    backgroundColor: '#222222',
                }}>
                    {displayItem(dataItemState)}
                </div>
                <div>
                    <WidgetItemHistory DataItemName={getDataItemName()} FullWidth={props.FullWidth}/>
                </div>
            </div>

        )
    }

    const {enqueueSnackbar} = useSnackbar();

    return (
        <div>
            <Grid container alignItems="center"
                  style={{backgroundColor: "#222", borderRadius: "10px", padding: "5px"}}>
                <Tooltip title="Back" TransitionComponent={Zoom}
                         style={{border: "1px solid #333", marginRight: "5px"}}>
                    <IconButton onClick={() => {
                        window.history.back()
                    }} variant="outlined" color="primary"><ArrowBackIosOutlinedIcon fontSize="large"/></IconButton>
                </Tooltip>
                <Divider orientation="vertical" flexItem style={{marginLeft: "20px", marginRight: "20px"}}/>
                <DialogAddItemPublicChannel ItemName={getDataItemName()} OnSuccess={(channel) => {
                    enqueueSnackbar("Item added to channel" + channel, {variant: 'success'});
                }}/>
                <DialogWriteValue ItemName={getDataItemName()} OnSuccess={() => {
                    enqueueSnackbar("Value has been written", {variant: 'success'});
                }}/>
            </Grid>

            {displayItems()}

            <div>
                <WidgetError Message={messageError}/>
            </div>
        </div>
    );
}
