export class StatusResponse {
  public artistCount: number;
  public channelCount: number;
  public imageLinksCount: number;
  public programEntryCount: number;
  public tagCount: number;
  public tvShowCount: number;
  public version: string;
  public serverDateTime: Date;
  public currentlyUpdating: boolean;
  public tvChannels: TvShow[];
}

export class TvShow {
  public id: number;
  public adapterFamily: string;
  public channelKey: string;
  public technicalId: string;
  public name: string;
  public homePage: string;
}

export interface Channel {
  adapterFamily?: Channel.AdapterFamilyEnum;
  channelKey?: Channel.ChannelKeyEnum;
  homePage?: string;
  id?: number;
  name?: string;
  technicalId?: string;
}
export namespace Channel {
  export type AdapterFamilyEnum = 'ARD' | 'ZDF';
  export const AdapterFamilyEnum = {
    ARD: 'ARD' as AdapterFamilyEnum,
    ZDF: 'ZDF' as AdapterFamilyEnum
  };
  export type ChannelKeyEnum = 'ARD' | 'ZDF' | 'ZDF_INFO' | 'ZDF_NEO' | 'DREISAT' | 'ARTE' | 'BR' | 'HR' | 'MDR' | 'NDR' | 'RBB' | 'RADIO_BREMEN_TV' | 'SR' | 'SWR_BW' | 'SWR_RP' | 'WDR' | 'ALPHA' | 'TAGESSCHAU_24' | 'ONE' | 'KIKA' | 'PHOENIX';
  export const ChannelKeyEnum = {
    ARD: 'ARD' as ChannelKeyEnum,
    ZDF: 'ZDF' as ChannelKeyEnum,
    ZDFINFO: 'ZDF_INFO' as ChannelKeyEnum,
    ZDFNEO: 'ZDF_NEO' as ChannelKeyEnum,
    DREISAT: 'DREISAT' as ChannelKeyEnum,
    ARTE: 'ARTE' as ChannelKeyEnum,
    BR: 'BR' as ChannelKeyEnum,
    HR: 'HR' as ChannelKeyEnum,
    MDR: 'MDR' as ChannelKeyEnum,
    NDR: 'NDR' as ChannelKeyEnum,
    RBB: 'RBB' as ChannelKeyEnum,
    RADIOBREMENTV: 'RADIO_BREMEN_TV' as ChannelKeyEnum,
    SR: 'SR' as ChannelKeyEnum,
    SWRBW: 'SWR_BW' as ChannelKeyEnum,
    SWRRP: 'SWR_RP' as ChannelKeyEnum,
    WDR: 'WDR' as ChannelKeyEnum,
    ALPHA: 'ALPHA' as ChannelKeyEnum,
    TAGESSCHAU24: 'TAGESSCHAU_24' as ChannelKeyEnum,
    ONE: 'ONE' as ChannelKeyEnum,
    KIKA: 'KIKA' as ChannelKeyEnum,
    PHOENIX: 'PHOENIX' as ChannelKeyEnum
  }
}

export interface ProgramEntry {
  adapterFamily?: ProgramEntry.AdapterFamilyEnum;
  channel?: Channel;
  createdAt?: string;
  description?: string;
  durationInMinutes?: number;
  endDateTime?: string;
  homePage?: string;
  id?: number;
  imageLinks?: Array<ImageLink>;
  startDateTime?: string;
  tags?: Array<Tag>;
  technicalId?: string;
  title?: string;
  updatedAt?: string;
  url?: string;
}
export namespace ProgramEntry {
  export type AdapterFamilyEnum = 'ARD' | 'ZDF';
  export const AdapterFamilyEnum = {
    ARD: 'ARD' as AdapterFamilyEnum,
    ZDF: 'ZDF' as AdapterFamilyEnum
  }
}
export interface Tag {
  createdAt?: string;
  id?: number;
  tagName?: string;
}
export interface ImageLink {
  createdAt?: string;
  id?: number;
  url?: string;
}
