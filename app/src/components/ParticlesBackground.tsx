import { useEffect, useState } from "react";
import Particles, { ParticlesProvider } from "@tsparticles/react";
import { loadSlim } from "@tsparticles/slim";
import type { Engine, ISourceOptions } from "@tsparticles/engine";

const initParticles = async (engine: Engine): Promise<void> => {
    try {
        await loadSlim(engine);
    } catch (error) {
        console.error("Failed to initialize tsParticles", error);
        throw error;
    }
};

const particlesOptions: ISourceOptions = {
    background: {
        color: {
            value: "transparent",
        },
    },
    detectRetina: true,
    fpsLimit: 60,
    fullScreen: {
        enable: false,
    },
    particles: {
        color: {
            value: "#4385c5",
        },
        links: {
            color: "#4385c5",
            distance: 150,
            enable: true,
            opacity: 0.2,
            width: 1,
        },
        move: {
            enable: true,
            speed: 1.2,
        },
        number: {
            density: {
                enable: true,
                width: 900,
                height: 900,
            },
            value: 55,
        },
        opacity: {
            value: {
                min: 0.15,
                max: 0.35,
            },
        },
        shape: {
            type: "circle",
        },
        size: {
            value: {
                min: 1,
                max: 2.5,
            },
        },
    },
    pauseOnBlur: true,
    pauseOnOutsideViewport: true,
};

export function ParticlesBackground() {
    const [motionAllowed, setMotionAllowed] = useState(false);

    useEffect(() => {
        const mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
        const updateMotionPreference = () => setMotionAllowed(!mediaQuery.matches);

        updateMotionPreference();
        mediaQuery.addEventListener("change", updateMotionPreference);
        return () => mediaQuery.removeEventListener("change", updateMotionPreference);
    }, []);

    return (
        <div className="auth-particles-background">
            {motionAllowed && (
                <ParticlesProvider init={initParticles}>
                    <Particles
                        id="auth-particles"
                        options={particlesOptions}
                        style={{ height: "100%", width: "100%" }}
                    />
                </ParticlesProvider>
            )}
        </div>
    );
}

export default ParticlesBackground;
