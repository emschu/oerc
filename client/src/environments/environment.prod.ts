/*
 * oerc, alias oer-collector
 * Copyright (C) 2021-2023 emschu[aet]mailbox.org
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
import packageInfo from '../../package.json';

export const environment = {
  production: true,
  serverEndpoint: '/api/v2',
  apiRequestTimeoutInSecs: 60,
  locale: 'de',
  appVersion: packageInfo.version,
  defaultSettingShowDeprecatedEntries: false,
};
