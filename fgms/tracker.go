
package fgms


import(
	"github.com/FreeFlightSim/go-fgms/tracker"
)
//  Add a tracking server
//  int -1 for fail or SUCCESS
func (me *FG_SERVER) AddTracker(host string, port int, isTracked bool){

	me.IsTracked = isTracked
	me.Tracker = tracker.NewFG_Tracker(host, port, 0)
	
	/* TODO
	#ifndef NO_TRACKER_PORT
	#ifdef USE_TRACKER_PORT
	if ( m_Tracker )
	{
		delete m_Tracker;
	}
	m_Tracker = new FG_TRACKER(Port,Server,0);
	#else // !#ifdef USE_TRACKER_PORT
	if ( m_Tracker )
	{
		msgctl(m_ipcid,IPC_RMID,NULL);
		delete m_Tracker;
		m_Tracker = 0; // just deleted
	}
	printf("Establishing IPC\n");
	m_ipcid         = msgget(IPC_PRIVATE,IPCPERMS);
	if (m_ipcid <= 0)
	{
		perror("msgget getting ipc id failed");
		return -1;
	}
	m_Tracker = new FG_TRACKER(Port,Server,m_ipcid);
	#endif // #ifdef USE_TRACKER_PORT y/n
	#endif // NO_TRACKER_PORT
	return (SUCCESS);
	*/
} // FG_SERVER::AddTracker()



// Updates the remote tracker  ?
//func (me *FG_SERVER) UpdateTracker(Callsign string, Passwd string, Modelname string, Timestamp int64, messType int) int {
func (me *FG_SERVER) UpdateTracker( player *FG_Player, messType int) int {
	//#ifndef NO_TRACKER_PORT
	//char            TimeStr[100];
	//mT_PlayerListIt CurrentPlayer;
	//Point3D         PlayerPosGeod;
	
	//string          Aircraft;
	//var Aircraft string = player.Aircraft()
	//string          Message;
	var Message string
	
	//tm              *tm;
	//FG_TRACKER::m_MsgBuffer buf;

	//if ! m_IsTracked || strcmp(Callsign.c_str(),"mpdummy") == 0)
	//{
	//	return (1);
	//}
	
	// Creates the UTC time string
	//tm = gmtime (& Timestamp);
	//tm := Now()
	/*sprintf (
		TimeStr,
		"%04d-%02d-%02d %02d:%02d:%02d",
		tm->tm_year+1900,
		tm->tm_mon+1,
		tm->tm_mday,
		tm->tm_hour,
		tm->tm_min,
		tm->tm_sec
	); */
	TimeStr :=  "2013-12-25 11.22.33"
	// Edits the aircraft name string
	/* size_t Index = Modelname.rfind ("/");
	if (Index != string::npos)
	{
		Aircraft = Modelname.substr (Index + 1);
	}
	else
	{
		Aircraft = Modelname;
	}
	Index = Aircraft.find (".xml");
	if (Index != string::npos)
	{
		Aircraft.erase (Index);
	}
	*/
	// Creates the message
	if messType == tracker.CONNECT {
	
		Message  = "CONNECT "
		Message += player.Callsign
		Message += " "
		Message += player.Passwd
		Message += " "
		Message += player.Aircraft()
		Message += " "
		//Message += TimeStr
		Message += TimeStr
		// queue the message
		//sprintf (buf.mtext, "%s", Message.c_str());
		//buf.mtype = 1;
	//#ifdef USE_TRACKER_PORT
		//pthread_mutex_lock( &msg_mutex ); // acquire the lock
		//msg_queue.push_back(Message); // queue the message
		//pthread_cond_signal( &condition_var );  // wake up the worker
		//pthread_mutex_unlock( &msg_mutex ); // give up the lock
	//#else // !#ifdef USE_TRACKER_PORT
		//msgsnd (m_ipcid, &buf, strlen(buf.mtext), IPC_NOWAIT);
	//#endif // #ifdef USE_TRACKER_PORT y/n
	//#ifdef ADD_TRACKER_LOG
		//write_msg_log(Message.c_str(), Message.size(), (char *)"IN: "); // write message log
	//#endif // #ifdef ADD_TRACKER_LOG
		me.TrackerConnect++ // count a CONNECT message queued
		return 0 //(0);
	
	}else if messType == tracker.DISCONNECT {
		Message  = "DISCONNECT "
		Message += player.Callsign
		Message += " "
		Message += player.Passwd
		Message += " "
		Message += player.Aircraft()
		Message += " "
		Message += TimeStr
		// queue the message
		//sprintf (buf.mtext, "%s", Message.c_str());
		//buf.mtype = 1;
	//#ifdef USE_TRACKER_PORT
		//pthread_mutex_lock( &msg_mutex ); // acquire the lock
		//msg_queue.push_back(Message); // queue the message
		//pthread_cond_signal( &condition_var );  // wake up the worker
		//pthread_mutex_unlock( &msg_mutex ); // give up the lock
	//#else // !#ifdef USE_TRACKER_PORT
		//msgsnd (m_ipcid, &buf, strlen(buf.mtext), IPC_NOWAIT);
	//#endif // #ifdef USE_TRACKER_PORT y/n
	//#ifdef ADD_TRACKER_LOG
		//write_msg_log(Message.c_str(), Message.size(),(char *)"IN: "); // write message log
	//#endif // #ifdef ADD_TRACKER_LOG
		//m_TrackerDisconnect++; // count a DISCONNECT message queued
		me.TrackerDisconnect++ // count a DISCONNECT message queued
		return 0 //(0);
	}
	
	// We only arrive here if type!=CONNECT and !=DISCONNECT
	//CurrentPlayer = m_PlayerList.begin();
	//while (CurrentPlayer != m_PlayerList.end())
	//{
	for _, CurrentPlayer := range me.Players {
		if CurrentPlayer.IsLocal {
		
			//sgCartToGeod (CurrentPlayer->LastPos, PlayerPosGeod);
			PlayerPosGeod := SG_CartToGeod( CurrentPlayer.LastPos )
			
			Message =  "POSITION "
			Message += CurrentPlayer.Callsign
			Message += " "
			Message += CurrentPlayer.Passwd;
			Message += " "
			//Message += NumToStr (PlayerPosGeod[Lat], 6)+" " //lat
			//Message += NumToStr (PlayerPosGeod[Lon], 6)+" " //lon
			//Message += NumToStr (PlayerPosGeod[Alt], 6)+" " //alt
			Message += PlayerPosGeod.ToSpacedString() + " "
			Message += TimeStr
			// queue the message
			//sprintf(buf.mtext,"%s",Message.c_str());
			//buf.mtype=1;
			//#ifdef USE_TRACKER_PORT
			//pthread_mutex_lock( &msg_mutex ); // acquire the lock
			//msg_queue.push_back(Message); // queue the message
			//pthread_cond_signal( &condition_var );  // wake up the worker
			//pthread_mutex_unlock( &msg_mutex ); // give up the lock
			//#else // !#ifdef USE_TRACKER_PORT
			//msgsnd(m_ipcid,&buf,strlen(buf.mtext),IPC_NOWAIT);
			//#endif // #ifdef USE_TRACKER_PORT y/n
			//m_TrackerPostion++; // count a POSITION messge queued
		}
		//Message.erase(0);
		//CurrentPlayer++;
	} // while
	//#endif // !NO_TRACKER_PORT
	return 0 //(0);
} // UpdateTracker (...)

