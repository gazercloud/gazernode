import React, {useEffect, useState} from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Grid from "@material-ui/core/Grid";
import Request, {RequestFailed} from "../request";
import Typography from "@material-ui/core/Typography";
import {Tooltip} from "@material-ui/core";
import Zoom from "@material-ui/core/Zoom";
import IconButton from "@material-ui/core/IconButton";
import EditIcon from '@material-ui/icons/Edit';
import {useSnackbar} from "notistack";
import PersonIcon from "@material-ui/icons/Person";
import DialogAddUser from "../Dialogs/DialogAddUser";
import MoreHorizIcon from '@material-ui/icons/MoreHoriz';
import WidgetError from "../Widgets/WidgetError";

export default function PageUsers(props) {
    const [users, setUsers] = React.useState({})
    const [hoverItem, setHoverItem] = useState("")
    const [messageError, setMessageError] = React.useState("")

    const { enqueueSnackbar } = useSnackbar();

    const btnStyle = (key) => {
        let borderTop = '0px solid #333333'
        let backColor = '#1E1E1E'
        let backColorHover = '#222222'

        if (hoverItem === key) {
            return {
                padding: '20px',
                minWidth: '300px',
                maxWidth: '300px',
                minHeight: '130px',
                maxHeight: '130px',
                borderRadius: '10px',
                marginRight: '10px',
                marginBottom: '10px',
                borderTop: borderTop,
                backgroundColor: backColor,
            }
        } else {
            return {
                padding: '20px',
                minWidth: '300px',
                maxWidth: '300px',
                minHeight: '130px',
                maxHeight: '130px',
                borderRadius: '10px',
                marginRight: '10px',
                marginBottom: '10px',
                borderTop: borderTop,
                backgroundColor: backColorHover,
            }
        }
    }


    const btnClick = (name) => {
        const idHex = new Buffer(name).toString('hex');
        props.OnNavigate("#form=user&user="+ idHex)
    }

    const handleEnter = (ev, key) => {
        setHoverItem(ev)
    }

    const handleLeave = (ev, key) => {
        setHoverItem("")
    }

    const requestUsers = () => {
        let req = {
        }
        Request('user_list', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setUsers(result)
                            setMessageError("")
                        }
                    );
                    return
                }
                if (res.status === 500) {
                    res.json().then(
                        (result) => {
                            setMessageError(result.error)
                        }
                    );
                    return
                }
                RequestFailed()
            }).catch(res => {
            setMessageError(res.message)
            RequestFailed()
        });
    }

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        props.OnTitleUpdate('Users')
        requestUsers()
        setFirstRendering(false)
    }

    useEffect(() => {
        const timer = setInterval(() => {
            requestUsers()
        }, 2000);
        return () => clearInterval(timer);
    });

    const displayItem = (item) => {
        return (
            <Grid container direction='column'>
                <Grid item onClick={btnClick.bind(this, item)} style={{cursor: 'pointer'}}>
                    <Grid container>
                        <Grid item>
                            <Grid container alignItems='center'>
                                <PersonIcon color='primary' fontSize='large' style={{marginRight: '5px'}} />
                                <Typography style={{
                                    fontSize: '14pt',
                                    overflow: 'hidden',
                                    height: '20pt',
                                    width: '200px',
                                    textOverflow: 'ellipsis',
                                    whiteSpace: 'nowrap',
                                }}>{item}</Typography>
                            </Grid>
                        </Grid>
                    </Grid>
                </Grid>
                <Grid item style={{marginTop: '10px', backgroundColor: "#181818", borderRadius: "10px"}}>
                    <Grid container={true} direction="row" alignItems="center">
                        <Grid item style={{flexGrow: 1, paddingLeft: "10px"}}>
                        </Grid>
                        <Grid item>
                            <Tooltip title="Stop" TransitionComponent={Zoom}>
                                <IconButton variant='outlined' color='primary' style={{margin: '0px', opacity: "50%"}} disabled={item.status === 'stopped'}>
                                    <MoreHorizIcon fontSize="default" onClick={btnClick.bind(this, item)} />
                                </IconButton>
                            </Tooltip>
                        </Grid>
                    </Grid>

                </Grid>
            </Grid>
        )
    }

    return (
        <div>
            <Grid container alignItems="center" style={{backgroundColor: "#222", borderRadius: "10px", padding: "5px"}}>
                <DialogAddUser OnSuccess={()=> {
                    enqueueSnackbar("User created", {variant: "success"})
                    requestUsers()
                }} />
            </Grid>
            <Grid container direction="column" style={{marginTop: "20px"}}>
                <Grid item>
                    <Grid container direction="row">
                        {users !== undefined && users.items !== undefined ? users.items.map((item) => (
                            <Grid item
                                  key={"gazer-page_users-user-" + item}
                                  button
                                  onMouseEnter={() => handleEnter(item)}
                                  onMouseLeave={() => handleLeave(item)}
                                  style={btnStyle(item)}
                            >
                                {displayItem(item)}
                            </Grid>
                        )) : <div/>}
                    </Grid>
                </Grid>
            </Grid>
            <div>
                <WidgetError Message={messageError} />
            </div>
        </div>
    );
}
