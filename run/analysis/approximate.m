function approximate(grid)
  use('Interaction');

  filename = locate('approximate');
  surrogate = h5read(filename, '/surrogate');
  surrogate = surrogate.Surrogate;

  ni = double(surrogate.Inputs);
  no = double(surrogate.Outputs);
  indices = reshape(surrogate.Indices{1}, ni, [])';
  errors = abs(reshape(surrogate.Surpluses{1}, no, []))';
  errors(errors < 1e-10) = 0;

  Plot.figure(1200, 400);
  for i = 1:no
    subplot(no, 1, i);
    semilogy(errors(:, i));
    Plot.title('Output %d', i-1);
  end

  nn = size(indices, 1);

  nu = size(unique(indices, 'rows'), 1);
  if nu ~= nn
    warning('found %d nonunique indices out of %d', nn-nu, nn);
  end

  if nargin == 0, return; end

  nodes = zeros(nn, ni);
  nodes(:) = feval(grid, indices);

  switch ni
  case 2
    Plot.figure(600, 600);
    Plot.line(nodes(:, 1), nodes(:, 2), 'discrete', true, 'style', ...
      {'MarkerFaceColor', 'k', 'MarkerEdgeColor', 'None', 'MarkerSize', 3});
    Plot.title('Nodes');
    Plot.limit([0, 1], [0, 1]);
  end
end

function nodes = open(indices)
  nodes = zeros(size(indices));
  levels = levelize(indices);
  orders = orderize(indices);
  for i = 1:numel(indices)
    nodes(i) = double(orders(i)+1) / double(bitshift(uint64(2), levels(i)));
  end
end

function nodes = closed(indices)
  nodes = zeros(size(indices));
  levels = levelize(indices);
  orders = orderize(indices);
  for i = 1:numel(indices)
    if levels(i) == 0
      nodes(i) = 0.5;
    else
      nodes(i) = double(orders(i)) / double(bitshift(uint64(2), levels(i)-1));
    end
  end
end

function levels = levelize(indices)
  levels = bitshift(bitshift(indices, 58), -58);
end

function orders = orderize(indices)
  orders = bitshift(indices, -6);
end
