package org.emschu.oer.collector.reader;

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

import org.jsoup.Connection;
import org.jsoup.Jsoup;
import org.jsoup.nodes.Document;
import org.jsoup.select.Elements;
import org.emschu.oer.collector.reader.parser.ard.ProgramEntryParser;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.validation.constraints.Null;
import java.io.IOException;
import java.util.Map;
import java.util.logging.Logger;

/**
 * Use this class to fetch web pages or apis
 */
@Component
public class Fetcher {
    private static long counter = 0;

    private static String proxyHost;

    private static int proxyPort;

    protected static final Logger LOG = Logger.getLogger(Fetcher.class.getName());

    /**
     * you should prefer to use fetchDocument() instead
     *
     * @param url url to fetch
     * @return jsoup Document object
     */
    private static Document getDocument(String url, @Null Map<String, String> additionalHeaders) {
        boolean isArdLink = url.contains(ProgramEntryParser.ARD_HOST);

        try {
            Connection connection = Jsoup.connect(url)
                    .header("Accept", "text/html")
                    .header("Accept-Encoding", "gzip, deflate, br")
                    .header("Accept-Language", "de,en-US")
                    .userAgent("Mozilla/5.0 (Windows NT 6.1; rv:60.0) Gecko/20100101 Firefox/60.0")
                    .maxBodySize(0)
                    .timeout(20 * 1000);

            // apply additional headers
            if (additionalHeaders != null) {
                additionalHeaders.forEach((s, s2) -> connection.header(s, s2));
            }

            if (isArdLink) {
                connection.header("Host", "programm.ard.de");
            }
            if (proxyHost != null && proxyPort != 0) {
                counter++;
                LOG.finest("using proxy for connection. host:" + proxyHost + " with port: " + proxyPort);
                return connection.proxy(proxyHost, proxyPort).get();
            } else {
                counter++;
                return connection.get();
            }
        } catch (IOException e) {
            LOG.warning("io exception: " + e.getMessage());
            LOG.throwing(Fetcher.class.getName(), "getProgram", e);
        }
        return null;
    }

    /**
     * this method fetches the url until testSelection is present
     *
     * @param url url to fetch until a response is retrieved. max 5 retries
     * @param testSelection jsoup selector for detection of correct page fetch
     * @return jsoup Document object
     */
    public static Document fetchDocument(String url, String testSelection,@Null Map<String, String> headers) {
        if (url == null) {
            throw new IllegalArgumentException(String.format("invalid url '%s'", url));
        }
        Document document = null;
        int counter = 0;
        while(document == null) {
            document = Fetcher.getDocument(url, headers);
            if (counter > 4) {
                throw new IllegalStateException(String.format("could not fetch url '%s'", url));
            }
            counter++;
        }
        Elements testElement = document.select(testSelection);
        if (testElement == null || testElement.isEmpty()) {
            return fetchDocument(url, testSelection, null);
        }
        return document;
    }

    /**
     * delegating method for convenience for {@link #fetchDocument(String, String, Map)}
     *
     *
     * @param url
     * @param testSelection
     * @return
     */
    public static Document fetchDocument(String url, String testSelection) {
        return fetchDocument(url, testSelection, null);
    }


    public static long getCounter() {
        return counter;
    }

    protected static void increaseCounter() {
        counter++;
    }

    public static String getProxyHost() {
        return proxyHost;
    }

    @Value(value = "${oer.collector.proxy_host}")
    public void setProxyHost(String proxyHost) {
        Fetcher.proxyHost = proxyHost;
    }

    public static int getProxyPort() {
        return proxyPort;
    }

    @Value(value = "${oer.collector.proxy_port}")
    public void setProxyPort(String proxyPort) {
        try {
            Fetcher.proxyPort = Integer.parseInt(proxyPort);
        } catch (NumberFormatException nfe) {
            Fetcher.proxyPort = 0;
        }
    }
}
