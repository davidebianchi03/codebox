import { useEffect, useRef, useCallback } from 'react';
import { CodeboxNotification } from '../types/notifications';

interface UseNotificationsOptions {
    onNotification?: (notification: CodeboxNotification) => void;
    onError?: (error: string) => void;
    enabled?: boolean;
    maxRetries?: number;
    initialRetryDelay?: number;
}

export const useNotifications = ({
    onNotification,
    onError,
    enabled = true,
    maxRetries = 5,
    initialRetryDelay = 1000,
}: UseNotificationsOptions) => {
    const wsRef = useRef<WebSocket | null>(null);
    const isMountedRef = useRef(true);
    const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const retryCountRef = useRef(0);
    const retryDelayRef = useRef(initialRetryDelay);

    const getBackoffDelay = useCallback((retryCount: number): number => {
        const exponentialDelay = initialRetryDelay * Math.pow(2, retryCount);
        const maxDelay = 30000;
        const jitter = Math.random() * 1000;
        return Math.min(exponentialDelay + jitter, maxDelay);
    }, [initialRetryDelay]);

    const connect = useCallback(() => {
        if (!isMountedRef.current || !enabled) {
            return;
        }

        if (wsRef.current) {
            return;
        }

        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const serverUrl = import.meta.env.VITE_SERVER_URL || window.location.origin;
        const wsUrl = `${protocol}//${new URL(serverUrl).host}/api/v1/notifications`;

        try {
            wsRef.current = new WebSocket(wsUrl);

            wsRef.current.onopen = () => {
                if (!isMountedRef.current) return;
                retryCountRef.current = 0;
                retryDelayRef.current = initialRetryDelay;
            };

            wsRef.current.onmessage = (event) => {
                if (!isMountedRef.current) return;

                try {
                    const data = JSON.parse(event.data);
                    
                    if(data.type === 'ping') {
                        wsRef.current?.send(JSON.stringify({ type: 'pong' }));
                        return;
                    }

                    const notification: CodeboxNotification = data;
                    onNotification?.(notification);
                } catch (error) {
                    console.error('[Notifications] Error parsing message:', error);
                }
            };

            wsRef.current.onerror = (event) => {
                if (!isMountedRef.current) return;
                const errorMsg = 'WebSocket error in notifications hub';
                onError?.(errorMsg);
            };

            wsRef.current.onclose = () => {
                if (!isMountedRef.current) return;
                wsRef.current = null;

                if (retryCountRef.current < maxRetries) {
                    const delay = getBackoffDelay(retryCountRef.current);
                    retryCountRef.current++;
                    retryDelayRef.current = delay;

                    reconnectTimeoutRef.current = setTimeout(() => {
                        if (isMountedRef.current) {
                            connect();
                        }
                    }, delay);
                } else {
                    const errorMsg = `Failed to connect to notifications hub after ${maxRetries} attempts`;
                    console.error(`[Notifications] ${errorMsg}`);
                    onError?.(errorMsg);
                }
            };
        } catch (error) {
            console.error('[Notifications] Error creating WebSocket:', error);
            onError?.(String(error));
        }
    }, [enabled, onNotification, onError, maxRetries, initialRetryDelay, getBackoffDelay]);

    const disconnect = useCallback(() => {
        if (reconnectTimeoutRef.current) {
            clearTimeout(reconnectTimeoutRef.current);
            reconnectTimeoutRef.current = null;
        }

        if (wsRef.current) {
            wsRef.current.close();
            wsRef.current = null;
        }

        retryCountRef.current = 0;
        retryDelayRef.current = initialRetryDelay;
    }, [initialRetryDelay]);

    useEffect(() => {
        isMountedRef.current = true;

        if (enabled) {
            connect();
        }

        return () => {
            isMountedRef.current = false;
            disconnect();
        };
    }, [enabled, connect, disconnect]);

    return {
        disconnect,
        isConnected: wsRef.current?.readyState === WebSocket.OPEN,
        retryCount: retryCountRef.current,
    };
};
