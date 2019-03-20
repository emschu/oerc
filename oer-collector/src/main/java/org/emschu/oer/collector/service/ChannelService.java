package org.emschu.oer.collector.service;

/*-
 * #%L
 * oer-server
 * %%
 * Copyright (C) 2019 emschu[aet]mailbox.org
 * %%
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 * 
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 * 
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 * #L%
 */

import org.emschu.oer.core.model.Channel;
import org.emschu.oer.core.model.repository.ChannelRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.logging.Logger;

@Service
public class ChannelService {

    private static final Logger LOG = Logger.getLogger(ChannelService.class.getName());

    @Autowired
    private ChannelRepository channelRepository;

    public void init() {
        LOG.info("Init sender list");
        for (Channel.AdapterFamily family : getActivatedSenderFamilies()) {
            List<Channel> channelList = getRegisteredSendersOfFamily(family);
            LOG.info("Enabling " + channelList.size() + " channels of channel family: " + family.toString());
            LOG.finest(channelList.toString());
            channelList.forEach(this::addChannel);
        }
    }

    /**
     * This method provides senders this application can handle
     *
     * @param adapterFamily
     * @return
     */
    public List<Channel> getRegisteredSendersOfFamily(Channel.AdapterFamily adapterFamily) {
        List<Channel> channelList = new ArrayList<>();
        switch (adapterFamily) {
            case ARD:
                channelList.add(new Channel(Channel.AdapterFamily.ARD, Channel.ChannelKey.ARD, "28106", "ARD - Das Erste", "https://ard.de"));
                channelList.add(new Channel(Channel.AdapterFamily.ARD, Channel.ChannelKey.BR, "28107", "BR Fernsehen", "https://www.br.de/fernsehen/index.html"));
                channelList.add(new Channel(Channel.AdapterFamily.ARD, Channel.ChannelKey.HR, "28108", "HR Fernsehen", "https://www.hr-fernsehen.de/"));
                channelList.add(new Channel(Channel.AdapterFamily.ARD, Channel.ChannelKey.MDR, "28229", "MDR Fernsehen", "https://www.mdr.de/tv/"));
                channelList.add(new Channel(Channel.AdapterFamily.ARD, Channel.ChannelKey.NDR, "28226", "NDR Fernsehen", "https://www.ndr.de"));
                channelList.add(new Channel(Channel.AdapterFamily.ARD, Channel.ChannelKey.RBB, "28205", "RBB Fernsehen", "https://www.rbb-online.de/fernsehen/"));
                channelList.add(new Channel(Channel.AdapterFamily.ARD, Channel.ChannelKey.RADIO_BREMEN_TV, "28385", "Radio Bremen TV", "https://www.radiobremen.de/fernsehen/"));
                channelList.add(new Channel(Channel.AdapterFamily.ARD, Channel.ChannelKey.SR, "28486", "SR Fernsehen", "https://www.sr.de/sr/home/fernsehen/"));
                channelList.add(new Channel(Channel.AdapterFamily.ARD, Channel.ChannelKey.SWR_BW, "28113", "SWR BW Fernsehen", "https://www.swrfernsehen.de/tv-programm/"));
                channelList.add(new Channel(Channel.AdapterFamily.ARD, Channel.ChannelKey.SWR_RP, "28231", "SWR RP Fernsehen", "https://www.swrfernsehen.de/tv-programm/"));
                channelList.add(new Channel(Channel.AdapterFamily.ARD, Channel.ChannelKey.WDR, "28111", "WDR Fernsehen", "http://www.wdr.de/tv/"));
                channelList.add(new Channel(Channel.AdapterFamily.ARD, Channel.ChannelKey.ALPHA, "28487", "ARD ALPHA", "http://www.br.de/fernsehen/ard-alpha/"));
                channelList.add(new Channel(Channel.AdapterFamily.ARD, Channel.ChannelKey.TAGESSCHAU_24, "28721", "Tagesschau24", "http://programm.tagesschau24.de/"));
                channelList.add(new Channel(Channel.AdapterFamily.ARD, Channel.ChannelKey.ONE, "28722", "ARD One", "http://www.one.ard.de/"));
                channelList.add(new Channel(Channel.AdapterFamily.ARD, Channel.ChannelKey.KIKA, "28008", "KIKA", "http://www.kika.de/"));
                break;
            case ZDF:
                channelList.add(new Channel(Channel.AdapterFamily.ZDF, Channel.ChannelKey.ZDF, "zdf", "ZDF", "http://www.zdf.de/"));
                channelList.add(new Channel(Channel.AdapterFamily.ZDF, Channel.ChannelKey.ZDF_INFO, "zdfinfo", "ZDFinfo", "https://www.zdf.de/dokumentation/zdfinfo-doku"));
                channelList.add(new Channel(Channel.AdapterFamily.ZDF, Channel.ChannelKey.ZDF_NEO, "zdfneo", "ZDFneo", "https://www.zdf.de/sender/zdfneo"));
                channelList.add(new Channel(Channel.AdapterFamily.ZDF, Channel.ChannelKey.PHOENIX, "phoenix", "Phoenix", "http://www.phoenix.de/"));
                channelList.add(new Channel(Channel.AdapterFamily.ZDF, Channel.ChannelKey.DREISAT, "3sat", "3Sat", "https://www.3sat.de"));
                channelList.add(new Channel(Channel.AdapterFamily.ZDF, Channel.ChannelKey.ARTE, "arte", "ARTE", "https://www.arte.tv/de"));
                break;
            default:
                throw new IllegalArgumentException("channel family '" + adapterFamily.toString() + "' is not yet implemented.");
        }
        return channelList;
    }

    /**
     * method to handle possible duplicates of channels
     * @param newChannel
     */
    public void addChannel(Channel newChannel) {
        if (channelRepository.existsByTechnicalId(newChannel.getTechnicalId())) {
            // do nothing
            return;
        }
        LOG.info("Adding new channel " + newChannel);
        channelRepository.save(newChannel);
    }

    public List<Channel> getAllSendersByFamily(Channel.AdapterFamily adapterFamily) {
        return channelRepository.findAllByAdapterFamily(adapterFamily);
    }

    /**
     * this place defines activated channel families
     *
     * @return
     */
    public Channel.AdapterFamily[] getActivatedSenderFamilies() {
        Channel.AdapterFamily[] families = new Channel.AdapterFamily[2];
        families[0] = Channel.AdapterFamily.ARD;
        families[1] = Channel.AdapterFamily.ZDF;
        return families;
    }

    public Channel.AdapterFamily[] getAllSenderFamilies() {
        return Channel.AdapterFamily.values();
    }

    public Channel getChannel(long id) {
        Optional<Channel> channelOptional = channelRepository.findById(id);
        return channelOptional.orElse(null);
    }
}
