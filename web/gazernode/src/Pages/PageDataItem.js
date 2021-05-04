import React, {useEffect, useState} from 'react';
import Request from "../request";
import Button from "@material-ui/core/Button";

function PageDataItem(props) {
    const [dataItemState, setDataItemState] = React.useState([])

    const btnStyle = (key) => {
        if (currentItem === key) {
            return {
                borderBottom: '1px solid #333333',
                cursor: "pointer",
                backgroundColor: "#222222",
            }
        } else {
            if (hoverItem === key) {
                return {
                    borderBottom: '1px solid #333333',
                    cursor: "pointer",
                    backgroundColor: "#222222"
                }
            } else {
                return {
                    borderBottom: '1px solid #333333',
                    cursor: "pointer",
                    backgroundColor: "#1E1E1E"
                }
            }
        }
    }

    const [currentItem, setCurrentItem] = useState("")
    const [hoverItem, setHoverItem] = useState("")

    const btnClick = (ev, name) => {
    }

    const handleEnter = (ev, key) => {
        setHoverItem(ev)
    }

    const handleLeave = (ev, key) => {
        setHoverItem("")
    }

    const requestDataItem = (dataItemName) => {
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

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        const dataItemName = new Buffer(props.DataItemName, 'hex').toString();
        requestDataItem(dataItemName)
        setFirstRendering(false)
    }

    useEffect(() => {
        const dataItemName = new Buffer(props.DataItemName, 'hex').toString();
        const timer = setInterval(() => {
            requestDataItem(dataItemName)
        }, 500);
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

    const displayItemName = (item, mainItem) => {
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

    if (dataItemState === undefined || dataItemState.value === undefined) {
        return (
            <div>
                loading ...
            </div>
        )
    }

    return (
        <div>
            <Button variant='outlined' color='primary' style={{minWidth: '100px', marginBottom: '20px'}}
                    onClick={()=>{
                        window.history.back()
                    }}
            >
                Back
            </Button>
            <div style={{
                padding: '20px',
                borderRadius: '10px',
                margin: '10px',
                backgroundColor: '#222222',
            }}>
                {displayItem(dataItemState)}
            </div>
        </div>
    );
}

export default PageDataItem;
