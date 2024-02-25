/* Do not change, this code is generated from Golang structs */


export enum DudesSprite {
    agent = "agent",
    cat = "cat",
    dude = "dude",
    girl = "girl",
    santa = "santa",
    sith = "sith",
    random = "random",
}
export interface DudesGrowRequest {
    channelId: string;
    userId: string;
    userName: string;
    userLogin: string;
    color: string;
}
export interface DudesUserSettings {
    dudeColor?: string;
    dudeSprite?: DudesSprite;
    userId: string;
    userName: string;
    userLogin: string;
}