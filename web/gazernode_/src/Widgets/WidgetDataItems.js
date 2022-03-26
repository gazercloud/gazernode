import React, {useEffect, useState} from 'react';
import Grid from "@material-ui/core/Grid";
import AddBoxOutlinedIcon from '@material-ui/icons/AddBoxOutlined';
import IndeterminateCheckBoxOutlinedIcon from '@material-ui/icons/IndeterminateCheckBoxOutlined';
import Request from "../request";

export default function WidgetDataItems(props) {
    const [treeContent, setTreeContent] = useState({
        "path": props.Root,
        "children": [
        ]
    })

    const [values, setValues] = useState({})
    const [lastRoot, setLastRoot] = useState("")

    const requestItems = (path) => {
        console.log("requestItems")
        console.log(path)

        let req = {
        }
        Request('data_item_list_all', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            let workCopy = JSON.parse(JSON.stringify(treeContent));
                            console.log("requestItems result")
                            console.log(result)
                            console.log(workCopy)
                            let node = findNodeByPath(workCopy, path);
                            if (node !== undefined) {
                                let children = []

                                let addedItems = {}
                                let itemsToAdd = []

                                for (let i = 0; i < result.items.length; i++)
                                {
                                    let origItem = result.items[i]
                                    if (origItem.name.includes(path + "/") || path === "") {
                                        let origName = origItem.name
                                        console.log("WDI origName1", origName, "path", path)
                                        if (path !== "")
                                            origName = origName.substr(path.length + 1)
                                        //origName = origName.replace(path + "/", "")
                                        console.log("WDI origName2", origName)
                                        let indexOfSlash = origName.indexOf("/")
                                        if (indexOfSlash > -1) {
                                            origName = origName.substr(0, indexOfSlash)
                                        }

                                        if (addedItems[origName] !== true && origName.length > 0)  {
                                            addedItems[origName] = true

                                            let pathOfItem = path + "/" + origName
                                            if (path === "") {
                                                pathOfItem = origName
                                            }

                                            let item = {
                                                "name": origName,
                                                "path": pathOfItem,
                                                "expanded": false,
                                                "children": []
                                            }

                                            let workCopy = JSON.parse(JSON.stringify(values));
                                            workCopy[origItem.path] = "---"
                                            setValues(workCopy)

                                            itemsToAdd.push(item)
                                        }
                                    }
                                }

                                itemsToAdd.sort(function (a, b) {
                                    if (a.name < b.name) {
                                        return -1;
                                    }
                                    if (a.name > b.name) {
                                        return 1;
                                    }
                                    return 0;
                                })

                                for (let i = 0; i < itemsToAdd.length; i++) {
                                    children.push(itemsToAdd[i])
                                }

                                node.children = children
                                node.expanded = true
                            } else {
                                console.log("node not found!!!!!!!!!!!!!!")
                            }
                            setTreeContent(workCopy)
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

    const formatDate = (date) => {

        var dd = date.getDate();
        if (dd < 10) dd = '0' + dd;

        var mm = date.getMonth() + 1;
        if (mm < 10) mm = '0' + mm;

        var yy = date.getFullYear() % 100;
        if (yy < 10) yy = '0' + yy;

        return dd + '.' + mm + '.' + yy;
    }

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

    const [requestValuesProcessing, setRequestValuesProcessing] = useState(false)
    const requestValues = (paths) => {
        if (requestValuesProcessing) {
            return
        }
        setRequestValuesProcessing(true)

        let req = {
            "items": paths
        }
        Request('data_item_list', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {

                            let workCopy = JSON.parse(JSON.stringify(values));
                            for (let i = 0; i < result.items.length; i++) {
                                let dt1 = new Date(result.items[i].value.t)
                                let dt2 = ""
                                if (dt1 !== undefined) {
                                    dt2 = formatTime(dt1)
                                }

                                if (result.items[i] !== undefined && result.items[i].value !== undefined) {
                                    workCopy[result.items[i].name] = {
                                        "value": result.items[i].value.v,
                                        "dt": dt2,
                                        "uom": result.items[i].value.u
                                    }
                                }

                            }
                            setValues(workCopy)
                            setRequestValuesProcessing(false)

                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                }
                setRequestValuesProcessing(false)
            })
            .catch((err) => {
                setRequestValuesProcessing(false)
                //setErrorMessage("Unknown error")
            })
    }

    const getValue = (path) => {
        let valueItem = values[path]
        if (valueItem === undefined)
            return ""
        return valueItem.value
    }
    const getUOM = (path) => {
        let valueItem = values[path]
        if (valueItem === undefined)
            return ""
        return valueItem.uom
    }

    const getDateTime = (path) => {
        let valueItem = values[path]
        if (valueItem === undefined)
            return ""
        return valueItem.dt
    }

    const [currentItem, setCurrentItem] = useState("")
    const [hoverItem, setHoverItem] = useState("")

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering || lastRoot !== props.Root) {
        console.log("firstRendering!!!!!!!!!!!!!!!!!!!!")
        console.log(props.Root)
        setLastRoot(props.Root)

        treeContent.path = props.Root
        setCurrentItem("")

        setTreeContent({
            "path": props.Root,
            "children": [
            ]
        })

        requestItems(props.Root)

        setFirstRendering(false)
    }

    useEffect(() => {
        const timer = setInterval(() => {
            let paths = []
            for (let [key, value] of Object.entries(values)) {
                paths.push(key)
            }
            requestValues(getAllPaths(treeContent))
        }, 500);
        return () => clearInterval(timer);
    });

    const findNodeByPath = (currentItem, path) => {
        if (path === "" || path === undefined || path === "/")
            return currentItem

        if (currentItem === undefined)
            return undefined;
        if (currentItem.path === path)
            return currentItem
        if (currentItem.children === undefined)
            return undefined
        for (let i = 0; i < currentItem.children.length; i++) {
            let ret = findNodeByPath(currentItem.children[i], path)
            if (ret)
                return ret
        }
        return undefined
    }

    const getAllPaths = (currentItem) => {
        let res = []
        if (currentItem === undefined)
            return res;
        res.push(currentItem.path)
        if (currentItem.children === undefined)
            return res
        for (let i = 0; i < currentItem.children.length; i++) {
            let ret = getAllPaths(currentItem.children[i])
            res = res.concat(ret)
        }
        return res
    }

    const expandNode = (path) => {
        let node = findNodeByPath(treeContent,path)
        if (node === undefined)
            return
        if (node.header === true)
            return
        if (node.expanded) {
            let workCopy = JSON.parse(JSON.stringify(treeContent));
            let node = findNodeByPath(workCopy, path);
            if (node !== undefined) {
                node.children = []
                node.expanded = false
            }
            setTreeContent(workCopy)
        } else {
            requestItems(path)
        }
    }

    const drawNodeChildren = (item) => {
        if (item === undefined) {
            return (<div/>)
        }
        if (item.children === undefined) {
            return (<div/>)
        }
        return (
            <div>
                {item.children.map((ch) => (
                <div>{drawNodeItem(ch)}</div>
                ))}
            </div>
        )
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
                    backgroundColor: "#FFFFFF"
                }
            }
        }
    }

    const btnClick = (ev, key) => {
        setCurrentItem(ev)
        props.OnDataItemSelected(ev)
    }

    const handleEnter = (ev, key) => {
        setHoverItem(ev)
    }

    const handleLeave = (ev, key) => {
        setHoverItem("")
    }

    const drawNodeItem = (item) => {
        return (
            <Grid
                container
                direction="column"
                style={{borderTop: "1px dotted #CCCCCC", cursor: "pointer"}}
            >
                <Grid item>
                    <Grid container
                          direction="row"
                          justify="space-between"
                          alignContent="flex-end"
                          onMouseEnter={() => handleEnter(item.path)}
                          onMouseLeave={() => handleLeave(item.path)}
                          style={btnStyle(item.path)}
                          onClick={btnClick.bind(this, item.path)}
                    >

                        <Grid item style={{lineHeight: "0px"}}>
                            <Grid container direction="row" alignItems="center">
                                <Grid item onClick={expandNode.bind(this, item.path)}>
                                    <div style={{verticalAlign: "middle"}}>
                                    {
                                        item.expanded === true ?
                                            <IndeterminateCheckBoxOutlinedIcon style={{color: "#888888", margin: "0px"}}/>
                                            :
                                            <AddBoxOutlinedIcon style={{color: "#888888", margin: "0px"}}/>
                                    }
                                    </div>
                                </Grid>
                                <Grid item style={{maxWidth: "400px"}}>
                                    {item.name}
                                </Grid>
                            </Grid>
                        </Grid>
                        <Grid item>
                            <Grid container alignItems="center">
                                <Grid item style={{paddingRight: "10px"}}><div style={{textAlign: "right", width: "200px", maxWidth: "200px"}}> {getValue(item.path)}</div></Grid>
                                <Grid item style={{paddingRight: "10px"}}><div style={{textAlign: "left", width: "75px", maxWidth: "75px", color: "#AAAAAA"}}> {getUOM(item.path)}</div></Grid>
                                <Grid item style={{paddingRight: "10px"}}><div style={{textAlign: "right", width: "75px", maxWidth: "75px", color: "#AAAAAA"}}> {getDateTime(item.path)}</div></Grid>
                            </Grid>
                        </Grid>
                    </Grid>
                </Grid>
                <Grid item>
                    <div style={{paddingLeft: "20px"}}>
                        {drawNodeChildren(item)}
                    </div>
                </Grid>
            </Grid>
        )
    }

    const drawHeaderItem = (item) => {
        return (
            <Grid
                container
                direction="column"
                style={{backgroundColor: "#DDDDDD", borderTop: "1px dotted #CCCCCC", cursor: "pointer"}}
            >
                <Grid item>
                    <Grid container direction="row" justify="space-between" alignContent="flex-end">
                        <Grid item>
                            <Grid container direction="row" alignItems="center">
                                <Grid item>
                                </Grid>
                                <Grid item style={{paddingLeft: "10px"}}>
                                    Name
                                </Grid>
                            </Grid>
                        </Grid>
                        <Grid item>
                            <Grid container alignItems="center">
                                <Grid item style={{paddingRight: "10px"}}><div style={{textAlign: "right", width: "200px", maxWidth: "200px"}}>Value</div></Grid>
                                <Grid item style={{paddingRight: "10px"}}><div style={{textAlign: "left", width: "75px", maxWidth: "75px"}}>UOM</div></Grid>
                                <Grid item style={{paddingRight: "10px"}}><div style={{textAlign: "right", width: "75px", maxWidth: "75px"}}>Time</div></Grid>
                            </Grid>
                        </Grid>
                    </Grid>
                </Grid>
            </Grid>
        )
    }


    const headerItem = {
        name: "Name",
    }

    return (
        <div>
            {
                drawHeaderItem(headerItem)
            }
            {treeContent.children !== undefined ? treeContent.children.map((item) => (
                drawNodeItem(item)
            )) : <div/>}
        </div>
    );
}
