/*
 * oerc, alias oer-collector
 * Copyright (C) 2021 emschu[aet]mailbox.org
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public
 * License along with this program.
 * If not, see <https://www.gnu.org/licenses/>.
 */
export interface StatusResponse {
  channel_family_count: number;
  channel_count: number;
  program_entry_count: number;
  tv_show_count: number;
  image_links_count: number;
  version: number;
  server_date_time: number;
  data_start_time: number;
  data_end_time: number;
  tv_channels: Channel[];
  tv_channel_families: ChannelFamily[];
}

export interface ChannelResponse {
  data: Channel[];
  size: number;
}

export interface ProgramResponse {
  from: string;
  to: string;
  channel_id: number;
  size: number;
  program_list: ProgramEntry[];
}

export interface LogEntryResponse {
  elements: LogEntry[];
  size: number;
  offset: number;
  limit: number;
}

export interface LogEntry {
  id: number;
  created_at: string;
  updated_at: string;
  message: string;
}

export interface Recommendation {
  id: number;
  created_at: number;
  program_entry: ProgramEntry;
  program_entry_id: number;
  channel_id: number;
  start_date_time: string;
  end_date_time: string;
  keywords: string;
}

export interface Channel {
  id: number;
  created_at: string;
  title: string;
  hash: string;
  url: string;
  technical_id: string;
  homepage: string;
  channel_family_id: number;
}

export interface Pong {
  date: string;
  data: string;
}

export enum AdapterFamilyEnum {
  ARD = 'ARD',
  ZDF = 'ZDF',
  ORF = 'ORF',
  SRF = 'SRF',
}

export enum ChannelKeyEnum {
  ARD = 'ARD',
  ZDF = 'ZDF',
  ZDF_INFO = 'ZDF_INFO',
  ZDF_NEO = 'ZDF_NEO',
  DREISAT = 'DREISAT',
  ARTE = 'ARTE',
  BR = 'BR',
  HR = 'HR',
  MDR = 'MDR',
  NDR = 'NDR',
  RBB = 'RBB',
  RADIO_BREMEN_TV = 'RADIO_BREMEN_TV',
  SR = 'SR',
  SWR_BW = 'SWR_BW',
  SWR_RP = 'SWR_RP',
  WDR = 'WDR',
  ALPHA = 'ALPHA',
  TAGESSCHAU_24 = 'TAGESSCHAU_24',
  ONE = 'ONE',
  KIKA = 'KIKA',
  PHOENIX = 'PHOENIX',
}

export interface ProgramEntry {
  id: number;
  created_at: Date;
  last_check: Date;
  start_date_time: Date;
  end_date_time: Date;
  tags: string;
  image_links?: ImageLink[];
  description: string;
  duration_in_minutes: number;
  homepage: string;
  url: string;
  title: string;
  hash: string;
  technical_id: string;
  channel_family_id: number;
  channel_id: number;
}

export interface ImageLink {
  createdAt?: string;
  id: number;
  url?: string;
}

export interface ChannelFamily {
  id: number;
  title: string;
}
