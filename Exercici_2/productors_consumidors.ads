package productors_consumidors is

   protected type Monitor_Productor_Consumidor is

      procedure iniciar_maquina(n_clients: Integer);
      entry Consumir (n_refresc_consumit: Integer; n_refrescs_a_consumir: Integer; nom_client: String);
      entry Reposar (n_reposador: Integer; fi_reposador: out boolean);
      entry fi_client(clients_restants: out Integer);

      private
         max_refrescs_maquina: integer;
         total_clients: integer := 0;
         clients_acabats: integer := 0;
         n_refrescs_maquina: integer := 0;

   end Monitor_Productor_Consumidor;

end productors_consumidors;
