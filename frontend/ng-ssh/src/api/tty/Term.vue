<template>
  <v-container fluid fill-width fill-height v-resize="resizeListener">
    <v-snackbar v-model="overlay">
      {{ message }}
      <v-btn color="blue" text @click="snackbar = false">Close</v-btn>
    </v-snackbar>
  </v-container>
</template>

<script>

import {Terminal} from "xterm";
import UTF8Decoder from '@/api/tty/utf8-coder'
import { ConnectionFactory } from "@/api/tty/websocket";
import {ttyTabStore} from '@/api/tty/tty_store'

import 'xterm/css/xterm.css';

export const protocols = ["webtty"];

const msgInputUnknown = '0';
const msgInput = '1';
const msgPing = '2';
const msgResizeTerminal = '3';

const msgUnknownOutput = '0';
const msgOutput = '1';
const msgPong = '2';
const msgSetWindowTitle = '3';
const msgSetPreferences = '4';
const msgSetReconnect = '5';

export default {
    props: ["ttyConf"],
    data() {
        return {
            term: new Terminal({}),
            decoder: new UTF8Decoder(),
            connectionFactory: new ConnectionFactory(this.ttyConf.ttyServer, protocols),
            overlay: false,
            message: '',
            messageTimeout: 3000,
            messageTimer: null,
        }
    },
    mounted() {
        this.term.open(this.$el)
        
        this.open()
    },
    destroyed() {
        console.log('Vue Terminal Close...')
        this.connection.close()
    },
    watch: {
        geometry: function(dims) {
       
            if(ttyTabStore.getters.getCurrTtyId !== this.ttyConf.id) {
                // console.log('watch geometry change: ',dims)
                // TODO: Remove reliance on private API
                const core = this.term._core;

                // Force a full render
                if (this.term.rows !== dims.rows || this.term.cols !== dims.cols) {
                    // core._renderCoordinator.clear();
                    this.term.resize(dims.cols, dims.rows);
                }
            } 
        }
    },
    computed: {
        geometry: function() {
            return ttyTabStore.getters.getGeometry;
        }
    },
    methods: {
        proposeDimensions() {
            const core = this.term._core;
            const parentElementStyle = window.getComputedStyle(this.term.element.parentElement);
            const parentElementHeight = parseInt(parentElementStyle.getPropertyValue('height'));
            const parentElementWidth = Math.max(0, parseInt(parentElementStyle.getPropertyValue('width')));
            const elementStyle = window.getComputedStyle(this.term.element);
            const elementPadding = {
                top: parseInt(elementStyle.getPropertyValue('padding-top')),
                bottom: parseInt(elementStyle.getPropertyValue('padding-bottom')),
                right: parseInt(elementStyle.getPropertyValue('padding-right')),
                left: parseInt(elementStyle.getPropertyValue('padding-left'))
            };
            const elementPaddingVer = elementPadding.top + elementPadding.bottom;
            const elementPaddingHor = elementPadding.right + elementPadding.left;
            const availableHeight = parentElementHeight - elementPaddingVer;
            const availableWidth = parentElementWidth - elementPaddingHor - core.viewport.scrollBarWidth;
            const geometry = {
            cols: Math.floor(availableWidth / core._renderService.dimensions.actualCellWidth),
            rows: Math.floor(availableHeight / core._renderService.dimensions.actualCellHeight)
            };
            return geometry;
        },
        fit() {
            const dims = this.proposeDimensions();
       
            if(ttyTabStore.getters.getCurrTtyId === this.ttyConf.id) {
                ttyTabStore.dispatch('setGeometry', dims)
            } else {
                return
            }
            // TODO: Remove reliance on private API
            const core = this.term._core;

            // Force a full render
            if (this.term.rows !== dims.rows || this.term.cols !== dims.cols) {
                this.term.resize(dims.cols, dims.rows);
            }
        },
        info(){
            return { columns: this.term.cols, rows: this.term.rows };
        },
        resizeListener() {
            // console.log('resizeListern:',this.ttyConf['id'],'x:', window.innerWidth, 'y:', window.innerHeight)
            if(!this.term._core._renderService) {
                return
            }

            this.fit()
            this.term.scrollToBottom();
            this.showMessage(String(this.term.cols) + "x" + String(this.term.rows), this.messageTimeout);

            this.term.focus()
        },
        onInput(callback) {
            this.term.onData((data) => {
                callback(data);
            });
        },
        output(data) {
            this.term.write(this.decoder.decode(data));
        },
        showMessage(message, timeout) {
            
            this.message = message
            this.overlay = true

            if (this.messageTimer) {
                clearTimeout(this.messageTimer);

                this.messageTimer = null;
            }
            if (timeout > 0) {
                /*eslint consistent-this: [2, "that"]*/
                const that = this
                this.messageTimer = window.setTimeout(() => {
                    that.overlay = false
                }, timeout);
            }
            
        },
        setWindowTitle(title) {
            console.log(title)
        },
        setPreferences(value) {
            console.log(value)
        },
        deactivate() {
            // this.term.off("data");
            // this.term.off("resize");
            this.term.blur();
        },
        reset() {
            this.removeMessage();
            this.term.clear();
        },
        close() {
            window.removeEventListener("resize", this.resizeListener);
            this.term.destroy();
        },
        open: function() {
            let connection = this.connectionFactory.create();
            let pingTimer = null ;
            let reconnectTimeout = null;
            this.connection = connection;

            const resizeHandler = (colmuns, rows) => {
                connection.send(
                    msgResizeTerminal + JSON.stringify(
                        {
                            columns: colmuns,
                            rows: rows
                        }
                    )
                );
            };

            this.term.onResize((data) => {
                resizeHandler(data.cols, data.rows);
            })

            const setup = () => {
                connection.onOpen(() => {
                    const termInfo = this.info();

                    // console.log("args: ", args)
                    connection.send(JSON.stringify(
                        {
                            name: this.ttyConf.name,
                            'auth_token': this.authToken,
                        }
                    ));

                    this.resizeListener()
                    
                    this.onInput(
                        (input) => {
                            connection.send(msgInput + input);
                        }
                    );

                    pingTimer = window.setInterval(() => {
                        connection.send(msgPing)
                    }, 30 * 1000);
                });

                connection.onReceive((data) => {
                    const payload = data.slice(1);
                    switch (data[0]) {
                        case msgOutput:
                            this.output(atob(payload));
                            break;
                        case msgPong:
                            break;
                        case msgSetWindowTitle:
                            this.setWindowTitle(payload);
                            break;
                        case msgSetPreferences: {
                            const preferences = JSON.parse( payload );
                            this.setPreferences(preferences);
                            break;
                        }
                        case msgSetReconnect: {
                            const autoReconnect = JSON.parse(payload);
                            console.log("Enabling reconnect: " + autoReconnect + " seconds")
                            this.reconnect = autoReconnect;
                            break;
                        }
                    }
                });

                connection.onClose(() => {
                    clearInterval(pingTimer);
                    this.deactivate();
                    this.showMessage("Connection Closed", 0);
                    this.term.write('\x1B[1;3;31m Closed by tty server \x1B')
                    if (this.reconnect > 0) {
                        reconnectTimeout = window.setTimeout(() => {
                            connection = this.connectionFactory.create();
                            this.reset();
                            setup();
                        }, this.reconnect * 1000);
                    }
                });

                connection.open();
            }

            setup();
            return () => {
                clearTimeout(reconnectTimeout);
                connection.close();
            }
        }
    }
}
</script>